FROM golang:1.21.1-bullseye as base

ARG GOOS=linux
ARG GOARCH=amd64
ARG EXPORTER_NAME="blackbox-arp-exporter"

COPY go.mod /opt/src/
COPY go.sum /opt/src/

WORKDIR /opt/src/
RUN go mod download

COPY . /opt/src
RUN go build -o ${EXPORTER_NAME}

FROM debian:stretch-slim
COPY --from=base /opt/src/blackbox-arp-exporter /usr/local/bin/