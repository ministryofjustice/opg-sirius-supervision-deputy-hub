.PHONY: cypress

all: go-lint build-all unit-test scan pa11y lighthouse cypress down

lint: go-lint npm-lint

build:
	docker compose build --no-cache --parallel deputy-hub

build-all:
	docker compose build --parallel deputy-hub cypress test-runner json-server

npm-lint:
	docker compose run --rm npm
	docker compose run --rm npm run lint

npm-prettier:
	docker compose run --rm npm
	docker compose run --rm npm run prettier

go-lint:
	docker compose run --rm go-lint

gosec: setup-directories
	docker compose run --rm gosec

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots

setup-directories: test-results

unit-test: setup-directories
	docker compose run --rm test-runner

cypress: setup-directories build-all
	docker compose up -d --wait deputy-hub json-server
	docker compose run --rm cypress run --env grepUntagged=true

cypress-single: setup-directories
	docker compose up -d --wait deputy-hub json-server
	docker compose run --rm cypress run --spec cypress/e2e/$(SPEC)

axe: setup-directories build-all
	docker compose up -d --wait deputy-hub
	docker compose run --rm cypress run --env grepTags="@axe"

up:
	docker compose up --build -d deputy-hub

dev-up:
	docker compose run --rm npm
	docker compose run --rm npm run build
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up --build npm deputy-hub json-server

down:
	docker compose down
