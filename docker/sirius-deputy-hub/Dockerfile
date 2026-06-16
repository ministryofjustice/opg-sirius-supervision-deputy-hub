FROM golang:1.26-alpine@sha256:7a3e50096189ad57c9f9f865e7e4aa8585ed1585248513dc5cda498e2f41812c AS base
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN update-ca-certificates

FROM node:24.14.0-alpine3.23@sha256:7fddd9ddeae8196abf4a3ef2de34e11f7b1a722119f91f28ddf1e99dcafdf114 AS asset-env

WORKDIR /app

COPY web/assets web/assets
COPY package.json .
COPY yarn.lock .

RUN yarn install --frozen-lockfile --ignore-scripts
RUN yarn build

### Development with hot reload and debugger
FROM base AS dev

WORKDIR /app

COPY --from=asset-env /app/web/static web/static

RUN go install github.com/air-verse/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

FROM base AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY internal internal
COPY web/template web/template
COPY main.go .
COPY main_test.go .

RUN go mod download

RUN go build -a -installsuffix cgo -o /go/bin/opg-sirius-supervision-deputy-hub

FROM alpine:latest@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4

WORKDIR /go/bin

RUN apk upgrade --no-cache busybox libcrypto3 libssl3 musl musl-utils zlib
RUN apk --no-cache add tzdata
ENV TZ="Europe/London"

COPY --from=build-env /go/bin/opg-sirius-supervision-deputy-hub opg-sirius-supervision-deputy-hub
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["./opg-sirius-supervision-deputy-hub"]
