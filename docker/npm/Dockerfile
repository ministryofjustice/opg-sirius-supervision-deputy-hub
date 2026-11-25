FROM node:24-alpine3.24@sha256:116b800c3754b440ee0aa3ed2f20f5bf3d99a554bfd9adeea09ea8017f6e85f5

RUN apk add --no-cache git

WORKDIR /home/node/app

COPY package.json package-lock.json .npmrc ./

RUN npm ci --ignore-scripts

ENTRYPOINT ["sh"]
