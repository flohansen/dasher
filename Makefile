setup:
	mockgen -package mocks -source internal/routes/routes.go -destination internal/mocks/routes.go

run:
	go run cmd/main.go

test:
	go test ./... -race