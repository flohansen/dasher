FROM node:lts-alpine AS builder
WORKDIR /usr/src/app
COPY ./web .
RUN yarn
RUN yarn build

FROM nginx:stable-alpine
COPY --from=builder /usr/src/app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
