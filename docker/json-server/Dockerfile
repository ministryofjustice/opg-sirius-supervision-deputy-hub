## build json-server for mock API
FROM node:24.21.0-alpine3.23@sha256:159fe64649038c30f8cc1ec4be3af3a6e93e3648678c31294e2c5058dbeb99f3 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

RUN apk add --no-cache curl

ENTRYPOINT ["yarn", "serve"]
