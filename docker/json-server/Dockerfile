## build json-server for mock API
FROM node:24.18.0-alpine3.23@sha256:595398b0081eacda8e1c4c5b97b76cd1020e4d58a8ebcb4843b9bca1e79e7436 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

RUN apk add --no-cache curl

ENTRYPOINT ["yarn", "serve"]
