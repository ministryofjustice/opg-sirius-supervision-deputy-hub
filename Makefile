all: lint test build-all scan pa11y lighthouse cypress down

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v

install-test-runner:
	curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v1.6.4/gotestsum_1.6.4_darwin_amd64.tar.gz" | tar -xz -C /usr/local/bin gotestsum

test: install-test-runner
	gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

build:
	docker-compose -f docker/docker-compose.ci.yml build deputy-hub

build-all:
	docker-compose -f docker/docker-compose.ci.yml build --parallel deputy-hub json-server puppeteer cypress

scan:
	docker run --rm  -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy sirius-deputy-hub:latest

pa11y:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="pa11y-ci" puppeteer

lighthouse:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="lhci collect assert" puppeteer

lighthouse-ci:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="lhci autorun" puppeteer

.PHONY: cypress
cypress:
	docker-compose -f docker/docker-compose.ci.yml run cypress

down:
	docker-compose -f docker/docker-compose.ci.yml down
