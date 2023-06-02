# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY import-service /app/service
RUN chmod +x /app/service

ENV SR_IMPORT_PORT=8080

EXPOSE 8080

ENTRYPOINT [ "./service" ]
