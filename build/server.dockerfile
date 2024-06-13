FROM golang:1.22-alpine AS builder
RUN apk update && apk add gcc musl-dev

WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM scratch

COPY --from=builder /usr/src/app/main /main

EXPOSE 3000
CMD [ "/main" ]
