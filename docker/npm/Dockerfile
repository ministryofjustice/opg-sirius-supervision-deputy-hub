FROM node:24-alpine3.24@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53

RUN apk add --no-cache git

WORKDIR /home/node/app

COPY package.json package-lock.json .npmrc ./

RUN npm ci --ignore-scripts

ENTRYPOINT ["sh"]
