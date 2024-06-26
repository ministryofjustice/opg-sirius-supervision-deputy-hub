# OPG SIRIUS SUPERVISION DEPUTY HUB

### Major dependencies

-   [Go](https://golang.org/) (>= 1.21)
-   [docker-compose](https://docs.docker.com/compose/install/) (>= 1.27.4)

#### Installing dependencies locally:

-   `yarn install`
-   `go mod download`

---

## Local development

The application ran through Docker can be accessed on `localhost:8888/supervision/deputies/1`.

To enable debugging and hot-reloading of Go and Java Script files:

`make dev-up`

If you are using VSCode, you can then attach a remote debugger on port `2345`. The same is also possible in Goland.
You will then be able to use breakpoints to stop and inspect the application.

Additionally, hot-reloading is provided by Air, so any changes to the Go code (including templates)
will rebuild and restart the application without requiring manually stopping and restarting the compose stack.

To run your changes in the context of your local sirius environment:

```
make build
# switch to opg-sirius repo
make dev-up
```

Note that if you are integrating with local Sirius instead of the mock server, you will need to use a valid deputy id in the url.

### Without docker

Alternatively to set it up not using Docker use below. This hosts it on `localhost:1234`

-   `yarn install && yarn build `
-   `go build main.go `
-   `./main `

---

## Run Cypress tests

```
make build cypress
```
## Run Cypress tests in UI
- `make up`
- `yarn && yarn cypress`

---

### Run the unit/functional tests

test sirius files: `yarn test-sirius`
test server files: `yarn test-server`
Run all Go tests: `make unit-test`

---

## Formatting

This project uses the standard Golang styleguide, and can be autoformatting by running `gofmt -s -w .`.
To run the go linter run `make go-lint`.

To format .gotmpl files and other assets, we use Prettier, which can be run using `yarn fmt`.
To run the JS linter run `make yarn-lint`

---

## Feature Flagging

Features can be flagged in the UI using the `feature_flagged` template function. Add the feature's name to the
comma-separated environment variable `FEATURES` for the environments you want it flagged for, and then call the
function with your feature name in the template. You can then do what you want with it, e.g. apply the `hide` CSS class.

---
## CI Tests

To run the entire build locally on your machine as it would be ran in CI run
```make```
