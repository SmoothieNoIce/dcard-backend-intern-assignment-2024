# Dcard Backend Intern Assignment 2024 server

### Usage

```
$ cp config.json.example config.json
$ vim config.json
$ make 
```

### Migration

####  create schema
```
$ cd command/create_schema
$ go run create_schema.go
```

####  create migration

```
$ migrate create -ext sql -dir $ database/migration create_cool_table
```

### Tools

#### migration commands

```
$ make migrate-up
$ make migrate-down
$ make migrate-version
```

#### migrate failed

https://github.com/golang-migrate/migrate/blob/master/FAQ.md#what-does-dirty-database-mean

### Swagger

#### url

http://127.0.0.1:8788/swagger/index.html#