## build json-server for mock API
FROM node:24.22.0-alpine3.23@sha256:e44c87c7dd212720eb00469aeac89504ca2be77b8f0d6dbffb30665d5fcc59e1 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN apk add --no-cache git
RUN npm ci --ignore-scripts

ENTRYPOINT ["npm", "run", "serve"]
