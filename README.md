# gotemplate — Cookiecutter-шаблон Go-сервиса

Шаблон Go-сервиса на Clean Architecture + DDD (gRPC / grpc-gateway, Postgres/pgx,
proto-генерация через Makefile). Используется для быстрого создания новых сервисов
без ручного переименования.

## Требования

- [Cookiecutter](https://cookiecutter.readthedocs.io/): `pipx install cookiecutter`
  (или `brew install cookiecutter`)
- Go (для сборки сгенерированного проекта)

## Генерация нового сервиса

```bash
# из GitHub
cookiecutter gh:mechta-market/gotemplate

# или из локальной копии
cookiecutter /path/to/gotemplate
```

Будут заданы **два вопроса**:

| Переменная | Значение по умолчанию | Что задаёт |
|---|---|---|
| `project_slug` | `myservice` | имя сервиса — Go-модуль, `ServiceName`, имя proto-пакета (`<slug>_v1`), docker-образ, имя БД |
| `module_owner` | `github.com/mechta-market` | префикс Go-модуля и owner для CI-образа |

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
cookiecutter . --no-input -o /tmp/out
cd /tmp/out/myservice && go build ./...
```
