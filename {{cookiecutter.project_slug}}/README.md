# {{ cookiecutter.project_slug }}

### DB dump:

```
pg_dump --no-owner -Fc -U postgres {{ cookiecutter.project_slug }} -f ./{{ cookiecutter.project_slug }}.custom
```

### DB restore:

```
dropdb -U postgres {{ cookiecutter.project_slug }}
createdb -U postgres {{ cookiecutter.project_slug }}
pg_restore --no-owner -d {{ cookiecutter.project_slug }} -U postgres ./{{ cookiecutter.project_slug }}.custom
```

### Install `migrate` command-tool:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### Create new migration:

```
migrate create -ext sql -dir migrations mg_name
```

### Apply migration:

```
migrate -path migrations -database "postgres://localhost:5432/db_name?sslmode=disable" up
```

### Proto generation:

```
make generate-proto
```

The proto toolchain lives in `bin/` (git-ignored) and is installed by `make proto-tools`, which
`generate-proto` runs first: `protoc` of the version pinned in the Makefile (`PROTOC_VERSION`) and
plugins of the versions from `go.mod`. Global `protoc` and plugins are not used. Installed versions:

```
make proto-tools-version
```
