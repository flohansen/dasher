FROM golang:1.22-alpine AS builder
RUN apk update && apk add gcc musl-dev

WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=1 go build -o main cmd/main.go

FROM alpine:latest
RUN apk add --no-cache libc6-compat

COPY --from=builder /usr/src/app/main /main
COPY --from=builder /usr/src/app/migrations /migrations

EXPOSE 3000
CMD [ "/main" ]
