
## Running the app

```bash
# run storages, queue
$ docker-compose up -d

# start api
go run api/main.go

# start consumers (can use many times)
go run consumer/main.go

# start test requests script
go run test-requests/main.ru
```
