FROM golang:1.22-alpine AS builder
RUN apk update && apk add gcc musl-dev
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install go.uber.org/mock/mockgen@latest

WORKDIR /usr/src/app
COPY . .
RUN make generate
RUN CGO_ENABLED=1 go build -o main cmd/main.go

FROM alpine:latest
RUN apk add --no-cache libc6-compat

COPY --from=builder /usr/src/app/main /main
COPY --from=builder /usr/src/app/migrations /migrations

EXPOSE 3000
CMD [ "/main" ]
