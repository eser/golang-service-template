# .RECIPEPREFIX := $(.RECIPEPREFIX)<space>
BINARY_NAME=service-cli
TESTCOVERAGE_THRESHOLD=0

.PHONY: init
init:
	brew install pre-commit
	brew install make
	pre-commit install
	go install github.com/air-verse/air@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/jandelgado/gcov2lcov@latest

.PHONY: dev
dev:
	air; if [ $$? -ne 0 ]; then go run ./cmd/${BINARY_NAME}/; fi

.PHONY: build
build:
	go build -o ./tmp/dist/${BINARY_NAME} ./cmd/${BINARY_NAME}/

.PHONY: multiarch-build
multiarch-build:
	GOARCH=amd64 GOOS=darwin go build -o ./tmp/dist/${BINARY_NAME}-darwin-amd64 ./cmd/${BINARY_NAME}
	GOARCH=amd64 GOOS=linux go build -o ./tmp/dist/${BINARY_NAME}-linux-amd64 ./cmd/${BINARY_NAME}
	GOARCH=amd64 GOOS=windows go build -o ./tmp/dist/${BINARY_NAME}-windows-amd64 ./cmd/${BINARY_NAME}
	GOARCH=arm64 GOOS=darwin go build -o ./tmp/dist/${BINARY_NAME}-darwin-arm64 ./cmd/${BINARY_NAME}
	GOARCH=arm64 GOOS=linux go build -o ./tmp/dist/${BINARY_NAME}-linux-arm64 ./cmd/${BINARY_NAME}
	GOARCH=arm64 GOOS=windows go build -o ./tmp/dist/${BINARY_NAME}-windows-arm64 ./cmd/${BINARY_NAME}

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	go clean

.PHONY: run
run: build
	./tmp/dist/${BINARY_NAME}

.PHONY: test-api
test-api:
	cd ./deployments/api/ && \
	bru run ./ --env development && \
	cd ../../

.PHONY: test
test:
	go test -failfast -count 1 ./...

.PHONY: test-cov
test-cov:
	go test -failfast -count 1 -coverpkg=./... -coverprofile=${TMPDIR}cov_profile.out ./...
	# `go env GOPATH`/bin/gcov2lcov -infile ${TMPDIR}cov_profile.out -outfile ./cov_profile.lcov

.PHONY: test-view-html
test-view-html:
	go tool cover -html ${TMPDIR}cov_profile.out -o ${TMPDIR}cov_profile.html
	open ${TMPDIR}cov_profile.html

.PHONY: test-ci
test-ci: test-cov
	$(eval ACTUAL_COVERAGE := $(shell go tool cover -func=${TMPDIR}cov_profile.out | grep total | grep -Eo '[0-9]+\.[0-9]+'))

	@echo "Quality Gate: checking test coverage is above threshold..."
	@echo "Threshold             : $(TESTCOVERAGE_THRESHOLD) %"
	@echo "Current test coverage : $(ACTUAL_COVERAGE) %"

	@if [ "$(shell echo "$(ACTUAL_COVERAGE) < $(TESTCOVERAGE_THRESHOLD)" | bc -l)" -eq 1 ]; then \
    echo "Current test coverage is below threshold. Please add more unit tests or adjust threshold to a lower value."; \
    echo "Failed"; \
    exit 1; \
  else \
    echo "OK"; \
  fi

.PHONY: dep
dep:
	go mod download
	go mod tidy

.PHONY: lint
lint:
	`go env GOPATH`/bin/golangci-lint run

.PHONY: container-start
container-start:
	docker compose --file ./deployments/compose.yml up --detach

.PHONY: container-start-prod
container-start-prod:
	docker compose --file ./deployments/compose.production.yml up --detach

.PHONY: container-rebuild
container-rebuild:
	docker compose --file ./deployments/compose.yml up --detach --build

.PHONY: container-rebuild-prod
container-rebuild-prod:
	docker compose --file ./deployments/compose.production.yml up --detach --build

.PHONY: container-restart
container-restart:
	docker compose --file ./deployments/compose.yml restart

.PHONY: container-restart-prod
container-restart-prod:
	docker compose --file ./deployments/compose.production.yml restart

.PHONY: container-stop
container-stop:
	docker compose --file ./deployments/compose.yml stop

.PHONY: container-stop-prod
container-stop-prod:
	docker compose --file ./deployments/compose.production.yml stop

.PHONY: container-destroy
container-destroy:
	docker compose --file ./deployments/compose.yml down

.PHONY: container-destroy-prod
container-destroy-prod:
	docker compose --file ./deployments/compose.production.yml down

.PHONY: container-update
container-update:
	docker compose --file ./deployments/compose.yml pull

.PHONY: container-update-prod
container-update-prod:
	docker compose --file ./deployments/compose.production.yml pull

.PHONY: container-dev
container-dev:
	docker compose --file ./deployments/compose.yml watch

.PHONY: container-ps
container-ps:
	docker compose --file ./deployments/compose.yml ps --all

.PHONY: container-ps-prod
container-ps-prod:
	docker compose --file ./deployments/compose.production.yml ps --all

.PHONY: container-logs-all
container-logs-all:
	docker compose --file ./deployments/compose.yml logs

.PHONY: container-logs-all-prod
container-logs-all-prod:
	docker compose --file ./deployments/compose.production.yml logs

.PHONY: container-logs
container-logs:
	docker compose --file ./deployments/compose.yml logs go-service

.PHONY: container-logs-prod
container-logs-prod:
	docker compose --file ./deployments/compose.production.yml logs go-service

.PHONY: container-cli
container-cli:
	docker compose --file ./deployments/compose.yml exec go-service bash

.PHONY: container-cli-prod
container-cli-prod:
	docker compose --file ./deployments/compose.production.yml exec go-service bash

.PHONY: container-push
container-push:
ifdef VERSION
	docker build --platform=linux/amd64 -t acikyazilim.registry.cpln.io/golang-service-template:v$(VERSION) .
	docker push acikyazilim.registry.cpln.io/golang-service-template:v$(VERSION)
else
	@echo "VERSION is not set"
endif
