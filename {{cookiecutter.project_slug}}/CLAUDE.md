# CLAUDE.md

Руководство для Claude Code по работе с этим репозиторием. Go-сервис на Clean Architecture + DDD,
с gRPC/grpc-gateway, Postgres (pgx), proto-генерацией через Makefile.

Конвенции этого стека вынесены в глобальные Claude Code скиллы (`crud`, `mobone`,
`golang-service`, `golang-samber-lo`) — они подхватываются автоматически по описанию.

---

## Структура проекта

### Верхний уровень
- `cmd/main.go` — entrypoint, поднимает `internal/app.App`.
- `internal/` — бизнес-логика и инфраструктура (закрытые пакеты).
- `api/proto/` — исходные `.proto` данного сервиса.
- `pkg/proto/` — сгенерированный код protobuf/grpc/gateway (**не редактировать вручную**).
- `migrations/` — SQL миграции Postgres.
- `docs/` — swagger JSON и статические доки, выдаются через `/docs/*`.
- `vendor-proto/` — внешние `.proto` зависимости (обновляются Makefile).
- `Dockerfile`, `Makefile` — сборка, запуск, генерация proto.
- `.env.example`, `.migrate_scripts.example` — примеры окружения и миграций.

### Внутренние пакеты (`internal/`)
- `internal/app/` — сборка приложения: серверы, миграции, метрики, трассировка, HTTP-gateway, DI.
  - `app.go` — граф зависимостей и запуск компонентов.
  - `grpc.go` — gRPC сервер + интерсепторы ошибок/метрик/трейсинга.
  - `grpc_gateway.go` — HTTP-gateway, CORS, error handler.
  - `system_http_server.go` — системный HTTP-сервер (порт 3003): /healthcheck, /docs/*, /metrics,
    манифест для pulse.
  - `manifest.go` — сведения о сервисе для манифеста pulse (от человека: название, описание, команда,
    критичность).
  - `migration.go` — запуск миграций из `migrations/`.
- `internal/config/` — конфигурация через env (см. `config.go`).
- `internal/handler/` — транспортный слой.
  - `grpc/` — gRPC handlers.
  - `grpc/dto/` — преобразование protobuf ↔ domain models.
  - могут быть и другие транспортные каналы.
- `internal/usecase/` — usecase-слой (валидация, оркестрация сервисов и доменных сервисов).
- `internal/domain/` — доменная модель, сервисы и репозитории.
  - `*/model/` — доменные структуры (entity).
  - `*/service/` — доменные сервисы (инварианты/логика).
  - `*/repo/` — репозитории.
  - `common/` — общие модели/утилиты/PG базовый репозиторий.
- `internal/service/` — сервисы (фоновые/инфраструктурные), для переиспользования или выделения логики.
- `internal/infra/pulsekit/` — манифест сервиса для pulse (см. «Манифест для pulse»). Пакет — копия
  эталона из gotemplate: на месте не править, обновлять целиком из шаблона.
- `internal/errs/` и `internal/constant/` — общие коды ошибок и константы (`Version`, `Commit`,
  `BuiltAt` — из ldflags при сборке).

---

## Архитектура: слои и зависимости

- **Transport** (`internal/handler/grpc/*`):
  - Работает только с protobuf DTO и usecase-интерфейсами.
  - Не обращается напрямую к репозиториям и сервисам.
- **Usecase** (`internal/usecase/*`):
  - Входной слой от транспортного слоя (запросы от внешних систем).
  - Валидация входных параметров.
  - Оркестрация доменных сервисов и сервисов `internal/service/*`.
  - Вход/выход — доменные модели (не protobuf).
  - Желательно не обращается в соседние usecases.
- **Domain** (`internal/domain/*`):
  - `model/` — структуры данных, сущности (entity).
    - **Запрещены теги сериализации** (`json:"..."`, `yaml:"..."` и т.п.) на полях доменной
      модели — она не знает про транспорт и про формат хранения.
  - `service/` — доменные операции и инварианты. Может использовать только репозиторий.
  - `repo/` — доступ к хранилищам. Может содержать подпапки для разных типов хранилищ или моков.
- **Service** (Infrastructure/background/external integrations, `internal/service/*`):
  - Фоновые процессы, интеграции с внешними системами.
  - Выделенные или переиспользуемые логики.
  - Может использовать другие сервисы `internal/service/*` и доменные сервисы `internal/domain/*/service`.
  - Не обращается в usecase слой.
- **Composition** (`internal/app/`):
  - Сборка зависимостей, запуск серверов, миграций, фоновых сервисов.

### Правило зависимостей
```
handler        → usecase
usecase        → domain service
usecase        → service
service        → service
service        → domain service
domain service → repo
```
- Обратные зависимости **запрещены**.
- К `repo` слою доступ только из `domain service`.

### Хранилища
- Postgres: все доменные entities (см. `migrations/*`).
- **Имена таблиц всегда в единственном числе**, без plural: `usr`, `app`, `secret`, `item`
  (не `usrs`, `apps`, `secrets`, `items`). То же значение указывается в `TableName` репозитория.

### Миграции
- Файлы в `migrations/` в формате `NNNNNN_<name>.up.sql` / `.down.sql` (golang-migrate).
- В `down`-миграциях во всех командах `DROP` обязательно указывать `CASCADE`
  (напр. `drop table if exists <table> cascade;`).
- В `down` объекты удаляются в порядке, обратном `up` (с учётом внешних ключей).

### API
- gRPC сервисы: `api/proto/{{ cookiecutter.__proto_pkg }}/*`.
- HTTP-gateway: через grpc-gateway + swagger (`docs/api.swagger.json`).
- DTO маппинг: `internal/handler/grpc/dto`.
- **REST-пути (route paths) всегда в единственном числе**, без plural:
  `/secret`, `/secret/{id}`, `/app`, `/item` (не `/secrets`, `/apps`, `/items`).
- **Пагинация: во всех `list` запросах `page` начинается с `0`** (нумерация страниц с нуля).
- **Для update использовать HTTP-метод `PUT`** (не `PATCH`): `put: "/<entity>/{id}"`.

### Ошибки и валидация
- Семантические ошибки — через `internal/errs` (см. gRPC interceptor в `internal/app/grpc.go`).
- Валидация параметров — в usecase.
- Нельзя пробрасывать ошибки наружу без wrapping (оборачивать в `fmt.Errorf("...: %w")`).
- Для работы с ошибками всегда используй `errors.Is` и `errors.AsType`. Избегай прямого
  сравнения ошибок (`==`) и type assertion (`err.(*MyError)`), чтобы корректно обрабатывать
  обёрнутые ошибки.

### Правила изменения кода
- gRPC DTO не должны протекать в доменные сервисы.
- `pkg/proto` и `docs/api.swagger.json` — генерируемые файлы (обновляются через Makefile).
- В тестах всегда предпочитай `testify`: `require` для проверок, прерывающих тест,
  и `assert` для остальных утверждений.
- При реализации worker pool / параллельной обработки используй `errgroup`
  (golang.org/x/sync/errgroup), а не ручное управление горутинами через `sync.WaitGroup` + каналы.

---

## Композиционный корень (`internal/app/app.go`)

`app.go` — единственная точка композиции приложения: здесь собирается граф зависимостей
(`repo → service → usecase → handler`) и описывается жизненный цикл. Бизнес-логики тут нет —
только связывание компонентов и управление их запуском/остановкой.

### Тип `App`
- В поля выносится **только то, чем нужно управлять после `Init`**: то, что надо явно
  останавливать (серверы), ждать (фоновые сервисы/обработчики), закрывать (трассировщик,
  pgx pool), а также корневой `ctx` с его `ctxCancel` и `exitCode`.
- Локальные звенья графа (repo, service, usecase, handler), которые нужны только для сборки и
  сразу передаются дальше, **не** выносятся в поля — это локальные переменные внутри `Init`.

### Импорты Композиционного корня
- Группируются блоками с пустой строкой между группами: стандартная библиотека → внешние
  зависимости → внутренние пакеты проекта.
- Внутренние пакеты-конструкторы импортируются с суффиксом-алиасом `P`, и в алиасе прописывается весь путь в camel-case
  (например, `internal/service/mdm` -> `serviceMdmP`, `internal/handler/grpc` -> `handlerGrpcP`, `internal/service/2gis/service` -> `service2gisServiceP`).

### Методы-фазы жизненного цикла
Фиксированный набор методов, каждый делает ровно одно:
- `Init` — создание и связывание всех зависимостей.
- `PreStartHook` — действия перед стартом.
- `Start` — запуск серверов и фоновых сервисов.
- `Listen` — блокировка до сигнала ОС (`SIGINT`/`SIGTERM`).
- `Stop` — отмена контекста и graceful-остановка серверов.
- `WaitJobs` — ожидание завершения фоновых задач.
- `Exit` — закрытие ресурсов и выход с `exitCode`.

### Стиль `Init`
- Сборка идёт **сверху вниз в порядке зависимостей**: инфраструктура (логгер, трассировка,
  pgx pool, кэш, миграции) → доменные блоки → сервисы → серверы.
- Каждый логический блок предваряется коротким комментарием-меткой в нижнем регистре
  (`// ord`, `// checkout`, `// grpc server`).
- Внутри блока цепочка строится единообразно:
  `repo := ...New(...)` → `service := ...New(repo)` → `usecase := ...New(service)` →
  `handler := ...New(usecase)`.
- Блоки, которым не нужны внешние переменные, оборачиваются в анонимный блок `{ ... }` —
  для ограничения области видимости и визуального разделения.

### Обработка ошибок при инициализации
- Используется хелпер `errCheck(err, msg)`: на этапе сборки любая ошибка фатальна
  (лог + `os.Exit(1)`). Ошибки из `Init` наверх не пробрасываются.
- Один заранее объявленный `var err error` переиспользуется по ходу `Init`.

### Парность Start / WaitJobs / Stop
- Для каждого фонового компонента, запускаемого в `Start()`, есть симметричный вызов в
  `WaitJobs()` (`.Wait()`) и/или в `Stop()`. **Порядок перечисления компонентов одинаков во
  всех трёх методах** — это упрощает чтение и сверку.

### Конфигурация и опциональные компоненты
- Все параметры берутся из единого глобального конфига (`config.Conf.*`) прямо в месте
  использования.
- Опциональные компоненты включаются по условию на конфиг (`if config.Conf.X != "" { ... }`),
  с фолбэком на in-memory/no-op реализацию через общий интерфейс.

---

## Runtime и конфигурация

### Запуск
- Entry: `cmd/main.go` → `internal/app.App`.
- На старте выполняются:
  - загрузка env (autoload `.env`),
  - настройка логгера/метрик/трейсинга,
  - pgx pool,
  - миграции (`internal/app/migration.go`),
  - запуск gRPC + HTTP-gateway.

### Переменные окружения
- Описаны в `internal/config/config.go`.
- Примеры: `.env.example`, `.migrate_scripts.example`.

### Системный HTTP-сервер
- Отдельный сервер на порту `3003` (`const systemHttpPort` в `internal/app/system_http_server.go`).
- Обслуживает служебные ручки: `/healthcheck`, `/docs/*`, `/metrics`.

### Манифест для pulse
- Сервис сам рассказывает pulse о себе на системном порту: `/.well-known/pulse` (манифест),
  `/.well-known/pulse/status` (состояние зависимостей), диагностические ручки `/diag/*`. Стандарт и
  правила конфиденциальности — `github.com/mechta-market/pulse`, `docs/service-manifest.md`;
  **агент, внедряющий манифест, сначала читает там раздел «Для агента»**.
- Зависимость объявляется одной строкой рядом с созданием клиента в `app.go`:
  `a.pulsekit.Depend("pg", "postgres", pulsekit.Host(config.Conf.PgDsn), true, a.pgpool.Ping)` — id,
  вид, хост, критичность, проверка. Адрес из конфигурации — только через `pulsekit.Host` (URL, DSN
  `host=… password=…`, `dns:///`, `host:port` → только хост). Postgres объявлен в шаблоне.
- Что ломает зависимость — `Depend(…).Affects("выдача заказов")`; платный API или API без
  health-ручки — пассивная проверка `pulsekit.NewPassive` (исходы настоящих вызовов, `Observe(err)`).
- Показатели состояния — `a.pulsekit.Gauge`/`GaugeTime` (глубина очереди, последняя выгрузка);
  свои метрики — `Metric`, узнаваемые ошибки в логах — `ErrorPattern`, инструкции — `Service.Runbooks`.
- Диагностическая ручка — `pulsekit.Handle(a.pulsekit, pulsekit.Endpoint{…}, fn)`: схема ответа
  строится из типа ответа (json-теги; тег `pulse:"personal=phone,maxLength=…,maxItems=…,enum=a|b,description=…"`,
  description — последним, запятые можно). Строковому параметру нужен `Pattern`, `Enum` или
  `Personal`; поля с именами секретов запрещены. Нарушение не роняет сервис: ручка не публикуется,
  причина — в лог и `kit.Problems()` (тест манифеста требует пустой список); длинные тексты —
  `kit.Warnings()`. Ошибка для человека — `pulsekit.Error{Status, Message}`,
  остальные ошибки наружу не уходят; своё сообщение из ошибки — `pulsekit.Describe(err)`.
- Вызов ручки pulsekit пишет в лог сам (с `X-Pulse-Request-Id`); для своих строк —
  `pulsekit.RequestId(ctx)`. Тест ручки — `kit.CheckEndpoint(id, params)`: ответ сверяется со схемой.
- Только чтение, быстро, без побочных эффектов; в ответ — только то, что можно показать на общем
  дашборде.

### Метрики и трассировка
- Prometheus метрики на `/metrics` (системный сервер, порт 3003) при `WITH_METRICS=true`.
  Используется `metrics.Registry`, а не дефолтный `promhttp.Handler()`.
- Трейсинг Jaeger включается при `WITH_TRACING=true` и `JAEGER_ADDRESS`.

### Документация и healthcheck
- `/healthcheck` — HTTP healthcheck (200 OK), порт 3003.
- `/docs/*` — статические docs + swagger (`docs/api.swagger.json`), порт 3003.

### Сборка
- `make build` создаёт бинарник `cmd/build/svc` (версия, полный коммит и время сборки — в ldflags:
  `build` манифеста).
- Dockerfile копирует бинарник, `docs/` и `migrations/` в `/app`.

### Flow проверки изменений
```
make generate-proto  →  gofmt  →  go test ./...  →  go run ./cmd/.
```
