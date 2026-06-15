FROM cypress/included:15.9.0@sha256:91111352c3adf4a0ff5bee6f3f49af56a66deb81d74a784dbb68bc98a3fbd26e

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe
RUN npm install "@cypress/grep@5.1.0"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
