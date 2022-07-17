# Majoo Test 2022
## Menjalankan Service
```bash
PSQL_HOST=<IP Database Server> PSQL_PORT=<Port Database Server> PSQL_DBNAME=<Database Name> PSQL_USER=<Database user name> PSQL_PASSWD=<Database User Password> API_HOST=<IP Web Service> API_PORT=<Port Web Service> go run github.com/epulskyyy/majoo-test-2022
```

## Menjalankan Test
```shell
go test -v ./... -coverprofile=cover.out  && go tool cover -html=cover.out
go test -v github.com/epulskyyy/majoo-test-2022/usecase -coverprofile=cover.out  && go tool cover -html=cover.out
```
