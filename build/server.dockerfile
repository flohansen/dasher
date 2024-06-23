FROM golang:1.22-alpine AS builder-go
RUN apk update && apk add gcc musl-dev

WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=1 go build -o main cmd/server/main.go

# ---

FROM node:lts-alpine AS builder-node

WORKDIR /usr/src/app
COPY ./web .
RUN yarn
RUN yarn build

# ---

FROM alpine:3.20
RUN apk add --no-cache libc6-compat

COPY --from=builder-go /usr/src/app/main /main
COPY --from=builder-go /usr/src/app/migrations /migrations
COPY --from=builder-node /usr/src/app/dist /dist

EXPOSE 3000
CMD [ "/main" ]
