# Dasher Server

![ci pipeline](https://github.com/flohansen/dasher-server/actions/workflows/main.yml/badge.svg)

## How to

### Run tests
```
make test
```
or
```
go test ./... -race
```
### Run server
```
make run
```
or
```
go run cmd/main.go
```

## API Paths
`GET` `/api/v1/features`: Get all features stored by the server \
`POST` `/api/v1/features`: Create or update feature