# syntax=docker/dockerfile:1

FROM golang:1.18.2-alpine3.16

WORKDIR /app

EXPOSE 8080

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /mesh-app

CMD [ "/mesh-app" ]