FROM cypress/included:15.20.1@sha256:327d62ee6d1c9bd7b4b490bb2772c1b70d37d56647953a4ec3105a939b7c61ff

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe
RUN npm install "@cypress/grep@5.1.0"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
