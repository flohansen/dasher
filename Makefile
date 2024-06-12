run:
	go run cmd/main.go

test:
	go test ./... -race

generate:
	sqlc generate
	mockgen -package mocks -source internal/routes/routes.go -destination internal/routes/mocks/routes.go