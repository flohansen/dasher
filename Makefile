run:
	go run cmd/main.go

test:
	go test ./... -race

gen:
	go generate ./...
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/feature.proto
