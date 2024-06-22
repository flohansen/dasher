run:
	go run cmd/main.go

test:
	go test ./... -race

generate:
	sqlc generate
	mockgen -package mocks -source internal/routes/routes.go -destination internal/routes/mocks/routes.go
	mockgen -package mocks -source internal/notification/feature.go -destination internal/notification/mocks/feature.go
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/feature.proto
