# gotemplate — Cookiecutter-шаблон Go-сервиса

Шаблон Go-сервиса на Clean Architecture + DDD (gRPC / grpc-gateway, Postgres/pgx,
proto-генерация через Makefile). Используется для быстрого создания новых сервисов
без ручного переименования.

## Требования

- [Cookiecutter](https://cookiecutter.readthedocs.io/): `pipx install cookiecutter`
  (или `brew install cookiecutter`)
- Go (для сборки сгенерированного проекта)
- `curl` и `unzip` — `make generate-proto` сам скачивает `protoc` и плагины закреплённых версий
  в `bin/` сервиса; ставить их глобально не нужно

## Генерация нового сервиса

из GitHub:
```bash
cookiecutter gh:rendau/gotemplate
```

или из локальной копии:
```bash
cookiecutter /path/to/gotemplate
```

Будут заданы вопросы:

| Переменная | Значение по умолчанию | Что задаёт |
|---|---|---|
| `project_slug` | `myservice` | имя сервиса — Go-модуль, `ServiceName`, имя proto-пакета (`<slug>_v1`), docker-образ, имя БД |
| `module_owner` | `github.com/rendau` | префикс Go-модуля и owner для CI-образа |
| `install_skills` | `y` | ставить ли скиллы Claude Code для этого стека — см. [Скиллы Claude Code](#скиллы-claude-code) |

Из них автоматически выводится всё остальное (производные переменные `_*`):

- Go-модуль: `<module_owner>/<project_slug>` (go.mod + все импорты)
- `ServiceName` (трейсинг/метрики) = `<project_slug>`
- Каталог и пакет proto = `api/proto/<project_slug>_v1`
- docker-образ (CI `IMAGE`) = `<project_slug>:latest`
- CI `OWNER` (`.github/workflows/deploy.yml`) = org из `module_owner`

Имя бинаря **статично** = `svc` (Makefile `BINARY_NAME`, Dockerfile `CMD`, `.gitignore`) —
одинаковое во всех проектах.

После генерации `post_gen` hook запускает `go mod tidy`, инициализирует git-репозиторий,
делает первый коммит (уже с итоговыми `go.mod`/`go.sum`) и, если на вопрос `install_skills` ответили «да», ставит скиллы
Claude Code. `go mod tidy` и установка скиллов — best-effort: если не получилось (нет сети,
нет `go`/`git`), сервис всё равно сгенерирован, а hook подсказывает, как доделать руками.

## Общие proto

Сервис использует общие `.proto`: `google/api` и `protoc-gen-openapiv2` (аннотации HTTP-ручек
и swagger) и `common/common.proto` (`ErrorRep`, `ListParamsSt`). Все они уже лежат в шаблоне
в `vendor-proto/`, так что сервис собирается без внешних репозиториев.

Обновлять их можно из общего репозитория `protos` организации — так `common.proto` остаётся
одинаковым у всех её сервисов, и они понимают ошибки друг друга. Договорённость такая: этот
репозиторий называется `protos` и лежит у владельца из `module_owner`. При генерации hook
проверяет SSH-доступ к `git@<хост>:<владелец>/protos.git` (например, для
`module_owner = github.com/acme` — `git@github.com:acme/protos.git`):

- доступ есть — адрес записывается в `PROTOS_REPO` в начале `Makefile` сервиса, и
  `make refresh-vendor-proto` обновляет `vendor-proto/` оттуда;
- репозитория нет или нет доступа — `PROTOS_REPO` пустой, `vendor-proto/` обновляется
  вместе с шаблоном. Указать репозиторий позже: `make refresh-vendor-proto PROTOS_REPO=<git-url>`
  или вписать его в `PROTOS_REPO`.

Проверка ничего не спрашивает (ни пароля, ни ключа) и не ждёт дольше нескольких секунд.

## Скиллы Claude Code

Конвенции стека (CRUD-сущность целиком, сервисный слой, repo-слой на mobone, `samber/lo`)
описаны в скиллах Claude Code из репозитория
[`rendau/go-service-skill`](https://github.com/rendau/go-service-skill): `crud`,
`golang-service`, а также смежные `mobone` и `golang-samber-lo`. Скиллы глобальные — ставятся
в `~/.claude/skills` один раз и работают во всех сервисах.

При `install_skills = y` hook:

- **скиллы уже установлены** — по симлинку `~/.claude/skills/crud` находит клон
  go-service-skill и предлагает его обновить: `git pull --ff-only` и повторный `install.sh`;
- **скиллов нет** — спрашивает, куда склонировать go-service-skill (по умолчанию —
  `go-service-skill` рядом с папкой нового сервиса; относительный путь считается от неё же),
  клонирует и запускает `install.sh`. Если по этому пути уже лежит клон — берёт его, занятую
  чужими файлами папку не трогает.

`install.sh` подключает скиллы симлинками на клон, поэтому дальше они обновляются через
`git pull` в клоне.

Без терминала (например, в CI) вопросы не задаются: установленные скиллы обновляются, а новый
клон кладётся в `GO_SERVICE_SKILL_DIR` или в путь по умолчанию. Отключить установку без вопроса:

```bash
cookiecutter gh:rendau/gotemplate --no-input install_skills=no ...
```

Поставить скиллы без генерации сервиса:

```bash
git clone https://github.com/rendau/go-service-skill.git
go-service-skill/install.sh
```

## Что внутри сгенерированного проекта

Минимальный рабочий скелет: `cmd/`, `internal/{app,config,constant,domain/common,errs,infra/metrics,util}`,
`api/proto/<slug>_v1/`, `pkg/proto/common`, `vendor-proto/`, `migrations/`, `docs/`,
Makefile, Dockerfile, CI и `.claude/settings.json`. Скиллов в самом проекте нет — они
глобальные, см. [Скиллы Claude Code](#скиллы-claude-code).

Интеграции с [pulse](https://github.com/rendau/pulse) (манифест сервиса для каталога и
диагностики) в шаблоне нет: она необязательна и делается отдельно, когда разработчик об этом
попросит, — по `docs/service-manifest.md` в репозитории pulse.

## Сопровождение шаблона

- Тело проекта живёт в каталоге `{{cookiecutter.project_slug}}/`.
- Переменные — в `cookiecutter.json`; хуки — в `hooks/`.
- Плейсхолдеры — синтаксис Jinja2 (`{{ "{{ cookiecutter.project_slug }}" }}` и т.п.).
- Сам шаблон-репозиторий **не компилируется** напрямую (в импортах — плейсхолдеры);
  это нормально. Корректность проверяется smoke-CI
  (`.github/workflows/template-ci.yml`): рендер с дефолтами + `go build`.

### Локальная проверка

```bash
cookiecutter . --no-input -o /tmp/out install_skills=no
cd /tmp/out/myservice && go build ./...
```
