FROM cypress/included:16.3.0@sha256:3efa8d156a094a6c13ddf3fba616e97f7492706685064b75a39823e1147eca97

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe
RUN npm install "@cypress/grep@7.0.0"

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
