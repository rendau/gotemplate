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

### Install grpc tools:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```
