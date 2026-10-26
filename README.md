# gotemplate — Cookiecutter-шаблон Go-сервиса

Шаблон Go-сервиса на Clean Architecture + DDD (gRPC / grpc-gateway, Postgres/pgx,
proto-генерация через Makefile). Используется для быстрого создания новых сервисов
без ручного переименования.

## Требования

- [Cookiecutter](https://cookiecutter.readthedocs.io/): `pipx install cookiecutter`
  (или `brew install cookiecutter`)
- Go (для сборки сгенерированного проекта)

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
| `service_title` | `<project_slug>` | человеческое название для каталога pulse |
| `service_description` | — (обязательно) | что делает сервис: словами, по которым его будут искать люди и агенты |
| `owner_team` | — (обязательно) | команда-владелец, одинаково во всех сервисах команды |
| `criticality` | `low` | `low` — внутреннее, `medium` — деградация заметна клиентам, `high` — без него не работают продажи/выдача |

Последние четыре — для манифеста сервиса (`internal/app/manifest.go`): по нему pulse узнаёт о
сервисе сам (стандарт — `github.com/rendau/pulse`, `docs/service-manifest.md`). Без команды и
описания генерация останавливается; кавычки в них не допускаются.

Из них автоматически выводится всё остальное (производные переменные `_*`):

- Go-модуль: `<module_owner>/<project_slug>` (go.mod + все импорты)
- `ServiceName` (трейсинг/метрики) = `<project_slug>`
- Каталог и пакет proto = `api/proto/<project_slug>_v1`
- docker-образ (CI `IMAGE`) = `<project_slug>:latest`
- CI `OWNER` (`.github/workflows/deploy.yml`) = org из `module_owner`

Имя бинаря **статично** = `svc` (Makefile `BINARY_NAME`, Dockerfile `CMD`, `.gitignore`) —
одинаковое во всех проектах.

После генерации `post_gen` hook инициализирует git-репозиторий, делает первый коммит
и (best-effort) запускает `go mod tidy`.

## Что внутри сгенерированного проекта

Минимальный рабочий скелет: `cmd/`, `internal/{app,config,constant,domain/common,errs,infra/metrics,util}`,
`api/proto/<slug>_v1/`, `pkg/proto/common`, `vendor-proto/`, `migrations/`, `docs/`,
Makefile, Dockerfile, CI, а также skill `crud` в `.claude/skills/` для генерации CRUD-сущностей.

## Сопровождение шаблона

- Тело проекта живёт в каталоге `{{cookiecutter.project_slug}}/`.
- Переменные — в `cookiecutter.json`; хуки — в `hooks/`.
- Плейсхолдеры — синтаксис Jinja2 (`{{ "{{ cookiecutter.project_slug }}" }}` и т.п.).
- Сам шаблон-репозиторий **не компилируется** напрямую (в импортах — плейсхолдеры);
  это нормально. Корректность проверяется smoke-CI
  (`.github/workflows/template-ci.yml`): рендер с дефолтами + `go build`.

### Локальная проверка

```bash
cookiecutter . --no-input -o /tmp/out owner_team=platform service_description="Smoke-проверка шаблона"
cd /tmp/out/myservice && go build ./...
```
