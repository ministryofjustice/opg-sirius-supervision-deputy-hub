# OPG SIRIUS SUPERVISION DEPUTY HUB

### Major dependencies

-   [Go](https://golang.org/) (>= 1.16)
-   [docker-compose](https://docs.docker.com/compose/install/) (>= 1.27.4)

#### Installing dependencies locally:

-   `yarn install`
-   `go mod download`

---

## Local development

The application ran through Docker can be accessed on `localhost:8888/supervision/deputies/1`.

To enable debugging and hot-reloading of Go files:

`docker-compose -f docker/docker-compose.dev.yml up --build`

If you are using VSCode, you can then attach a remote debugger on port `2345`. The same is also possible in Goland.
You will then be able to use breakpoints to stop and inspect the application.

Additionally, hot-reloading is provided by Air, so any changes to the Go code (including templates)
will rebuild and restart the application without requiring manually stopping and restarting the compose stack.

To develop with your local sirius environment:

`docker-compose -f docker/docker-compose.yml up -d --build`

Note that if you are integrating with local Sirius instead of the mock server, you will need to use a valid deputy id in the url.

### Without docker

Alternatively to set it up not using Docker use below. This hosts it on `localhost:1234`

-   `yarn install && yarn build `
-   `go build main.go `
-   `./main `

---

## Run Cypress tests

`docker-compose -f docker/docker-compose.dev.yml up -d --build `

`yarn && yarn cypress `

## Run Cypress tests for M1 chipset

`yarn cypress-for-m1-build` <br>
`yarn cypress-for-m1-up` <br>
`yarn cypress-headless` This command will run all the cypress tests <br>
`yarn cypress-headless --spec "cypress/integration/A_FILE_NAME.spec.js` will only run the tests in that file

Only thing to note is that if there are any changers locally to a gotmpl file you will need to run<br>
`cypress-build-down`<br>
Then the above commands to pull in the latest code to test against.<br>
You will not have to re-build anything tho if you change the code in a spec file only.

### Run the unit/functional tests

test sirius files: `yarn test-sirius`
test server files: `yarn test-server`
Run all Go tests: `go test ./...`

---

## Formatting

This project uses the standard Golang styleguide, and can be autoformatting by running `gofmt -s -w .`.

To format .gotmpl files and other assets, we use Prettier, which can be run using `yarn fmt`.
