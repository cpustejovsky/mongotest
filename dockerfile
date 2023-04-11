# syntax=docker/dockerfile:latest
FROM golang:1.19 as test

COPY . /mongotest

WORKDIR /mongotest/cmd/web

RUN go build -o mongotest

EXPOSE 80

EXPOSE 27017

CMD [ "./mongotest" ]