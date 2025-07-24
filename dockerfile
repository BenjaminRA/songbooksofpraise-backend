FROM golang:1.19-alpine

RUN apk update && apk add build-base

ENV GIN_MODE=release

WORKDIR /src/app

COPY . .

RUN go install

EXPOSE 8080

CMD himnario-backend