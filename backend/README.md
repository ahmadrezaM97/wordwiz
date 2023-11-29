
## Postgres

To run docker container:
```
docker run -d \
  --name wordwiz-postgres \
  -e POSTGRES_PASSWORD=wordwiz123qwerty \
  -p 5432:5432 postgres:latest
```

To create a database in psql:
```
CREATE USER wordwiz_user WITH PASSWORD 'wordwiz$123qwerty!';

CREATE DATABASE wordwiz_db;

GRANT ALL PRIVILEGES ON DATABASE wordwiz_db TO wordwiz_user;

```


## Migrate DB:
https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### create
```
migrate create -ext sql -dir internal/storage/postgres/migrations -seq -digits 2 init_tables
```

### up
```
migrate -path internal/storage/postgres/migrations -database 'postgres://wordwiz_user:wordwiz$123qwerty!@0.0.0.0:5432/wordwiz_db?sslmode=disable' up
```

### down
```
migrate -path internal/storage/postgres/migrations -database 'postgres://wordwiz_user:wordwiz$123qwerty!@0.0.0.0:5432/wordwiz_db?sslmode=disable' down
```

## swag init
```
swag init -g ./cmd/main.go -o ./internal/server/docs
```