# syntax=docker/dockerfile:1

FROM golang:1.20.5-alpine3.18

WORKDIR /app

COPY import-service /app/service
COPY config /app/config
COPY assets /app/assets
RUN chmod +x /app/service

ENV SR_IMPORT_PORT=8080

EXPOSE 8080

ENTRYPOINT [ "./service" ]
