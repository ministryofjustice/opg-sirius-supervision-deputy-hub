FROM --platform=linux/amd64 golang:1.23-alpine as base
WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

RUN update-ca-certificates

FROM node:20.2.0-alpine3.16 as asset-env

WORKDIR /app

COPY web/assets web/assets
COPY package.json .
COPY yarn.lock .

RUN yarn install
RUN yarn build

### Development with hot reload and debugger
FROM base AS dev

WORKDIR /app

COPY --from=asset-env /app/web/static web/static

RUN go install github.com/air-verse/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

FROM base as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o /go/bin/opg-sirius-supervision-deputy-hub

FROM alpine:latest

WORKDIR /go/bin

RUN apk upgrade --no-cache busybox libcrypto3 libssl3
RUN apk --no-cache add tzdata
ENV TZ="Europe/London"

COPY --from=build-env /go/bin/opg-sirius-supervision-deputy-hub opg-sirius-supervision-deputy-hub
COPY --from=build-env /app/web/template web/template
COPY --from=asset-env /app/web/static web/static
ENTRYPOINT ["./opg-sirius-supervision-deputy-hub"]
