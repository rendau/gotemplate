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
| `install_skills` | `y` | ставить ли скиллы Claude Code для этого стека — см. [Скиллы Claude Code](#скиллы-claude-code) |

`service_title`, `service_description`, `owner_team` и `criticality` — для манифеста сервиса (`internal/app/manifest.go`): по нему pulse узнаёт о
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

После генерации `post_gen` hook инициализирует git-репозиторий, делает первый коммит,
запускает `go mod tidy` и, если на вопрос `install_skills` ответили «да», ставит скиллы
Claude Code. `go mod tidy` и установка скиллов — best-effort: если не получилось (нет сети,
нет `go`/`git`), сервис всё равно сгенерирован, а hook подсказывает, как доделать руками.

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

## Сопровождение шаблона

- Тело проекта живёт в каталоге `{{cookiecutter.project_slug}}/`.
- Переменные — в `cookiecutter.json`; хуки — в `hooks/`.
- Плейсхолдеры — синтаксис Jinja2 (`{{ "{{ cookiecutter.project_slug }}" }}` и т.п.).
- Сам шаблон-репозиторий **не компилируется** напрямую (в импортах — плейсхолдеры);
  это нормально. Корректность проверяется smoke-CI
  (`.github/workflows/template-ci.yml`): рендер с дефолтами + `go build`.

### Локальная проверка

```bash
cookiecutter . --no-input -o /tmp/out owner_team=platform service_description="Smoke-проверка шаблона" install_skills=no
cd /tmp/out/myservice && go build ./...
```
