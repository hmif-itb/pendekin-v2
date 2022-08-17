FROM golang:1.19.0-alpine3.15 AS build-be

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
COPY ./main.go .
RUN go build -o app

FROM node:16-alpine3.15 AS build-fe

WORKDIR /app
COPY view /app
RUN npm install
RUN npm run build

FROM alpine:3.15

ARG user=hmif

WORKDIR /app
COPY --from=build-be /app/app .
COPY --from=build-fe /app/public view/public
EXPOSE 3000
RUN addgroup --gid 1001 -S ${user} && \
    adduser -G ${user} --disabled-password --uid 1001 ${user} && \
    mkdir -p /var/log/${user} && \
    chown ${user}:${user} /var/log/${user}
USER ${user}
CMD ["/app/app"]
