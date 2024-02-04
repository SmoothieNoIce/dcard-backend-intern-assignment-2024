# Dcard Backend Intern Assignment 2024 server

### How to run this project

```
make 
```

### Migration

####  create migration

```
migrate create -ext sql -dir database/migration create_cool_table
```

### Tools

#### migration commands

```
make migrate-up
make migrate-down
make migrate-version
```

#### migrate failed

https://github.com/golang-migrate/migrate/blob/master/FAQ.md#what-does-dirty-database-mean

### Swagger

#### url

http://127.0.0.1/swagger/index.html#