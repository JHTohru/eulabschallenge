# Products API

This application is a pretty simple Product CRUD REST API written in Go and which uses Postgres as database.

## docker-compose

The docker-compose defined in [pkg/postgres/infra/docker-compose.yaml](pkg/postgres/infra/docker-compose.yaml) is needed to run the server ([cmd/restsvr/main.go](cmd/restsvr/main.go)) and to run the integration tests from ([pkg/postgres](pkg/postgres)).

### Running docker-compose:

```console
cd pkg/postgres/infra && docker-compose up
```

The host system must be a Linux system because the docker-compose uses [tmpfs mount](https://docs.docker.com/storage/tmpfs/).

Once the Postgres database is running, run the migration files to structure the database.

```console
go run cmd/migrate/main.go --dir ./pkg/postgres/infra/migrations up
```

## Starting the server

```console
go run cmd/restsvr/main.go
```

## Postman

There is a Postman collection file located in [pkg/echohandler/infra/postman_collection.json](pkg/echohandler/infra/postman_collection.json). It may be imported in Postman to ease the API testing.
