FROM cypress/included:15.13.0@sha256:1048a3bac8a29e63159668a36d6dbbcf3a8a291d6f413c85edaf7e39bba9f3d9

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe
RUN npm install "@cypress/grep@5.1.0"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
