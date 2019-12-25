FROM golang:latest

COPY . /sample-admission-controller
WORKDIR /sample-admission-controller
RUN GOOS=linux GOARCH=amd64 go build -o admission-controller

FROM alpine:latest
COPY --from=0 /sample-admission-controller/admission-controller /usr/local/bin/
