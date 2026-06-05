## build json-server for mock API
FROM node:24.14.0-alpine3.23 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

RUN apk add --no-cache curl

ENTRYPOINT ["yarn", "serve"]
