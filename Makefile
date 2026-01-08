.PHONY: cypress

all: go-lint build-all unit-test scan pa11y lighthouse cypress down

lint: go-lint yarn-lint

build:
	docker compose build --no-cache --parallel deputy-hub

build-all:
	docker compose build --parallel deputy-hub cypress test-runner json-server

yarn-lint:
	docker compose run --rm yarn
	docker compose run --rm yarn lint

yarn-prettier:
	docker compose run --rm yarn
	docker compose run --rm yarn prettier . --write

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots .trivy-cache

setup-directories: test-results

unit-test: setup-directories
	docker compose run --rm test-runner

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-deputy-hub:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-deputy-hub:latest

cypress: setup-directories build-all
	docker compose up -d --wait deputy-hub json-server
	docker compose run --rm cypress run --env grepUntagged=true

cypress-single: setup-directories build-all
	docker compose up -d --wait deputy-hub json-server
	docker compose run --rm cypress run --spec cypress/e2e/$(SPEC)

axe: setup-directories build-all
	docker compose up -d --wait deputy-hub
	docker compose run --rm cypress run --env grepTags="@axe"

up:
	docker compose up --build -d deputy-hub

dev-up:
	docker compose run --rm yarn
	docker compose run --rm yarn build
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up --build yarn deputy-hub json-server

down:
	docker compose down
