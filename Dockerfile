# syntax=docker/dockerfile:1

FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache qpdf

COPY import-service /app/service
COPY config /app/config
COPY assets /app/assets
RUN chmod +x /app/service

ENV SR_IMPORT_PORT=8080

RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Berlin /etc/localtime

EXPOSE 8080

ENTRYPOINT [ "./service" ]
