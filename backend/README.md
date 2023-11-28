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