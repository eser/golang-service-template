.PHONY: dev build multiarch-build run clean test test-cov test-ci dep lint container-start container-start-prod container-rebuild container-rebuild-prod container-restart container-restart-prod container-stop container-stop-prod container-destroy container-destroy-prod container-update container-update-prod container-dev container-ps container-ps-prod container-logs-all container-logs-all-prod container-logs container-logs-prod container-cli container-cli-prod container-push
# .RECIPEPREFIX := $(.RECIPEPREFIX)<space>
BINARY_NAME=service-cli
TESTCOVERAGE_THRESHOLD=0

dev:
	# go run ./cmd/${BINARY_NAME}/
	air

build:
	go build -o ./dist/${BINARY_NAME} ./cmd/${BINARY_NAME}/

multiarch-build:
	GOARCH=amd64 GOOS=darwin go build -o ./dist/${BINARY_NAME}-darwin-amd64 ./cmd/${BINARY_NAME}
	GOARCH=amd64 GOOS=linux go build -o ./dist/${BINARY_NAME}-linux-amd64 ./cmd/${BINARY_NAME}
	GOARCH=amd64 GOOS=windows go build -o ./dist/${BINARY_NAME}-windows-amd64 ./cmd/${BINARY_NAME}
	GOARCH=arm64 GOOS=darwin go build -o ./dist/${BINARY_NAME}-darwin-arm64 ./cmd/${BINARY_NAME}
	GOARCH=arm64 GOOS=linux go build -o ./dist/${BINARY_NAME}-linux-arm64 ./cmd/${BINARY_NAME}
	GOARCH=arm64 GOOS=windows go build -o ./dist/${BINARY_NAME}-windows-arm64 ./cmd/${BINARY_NAME}

run: build
	./dist/${BINARY_NAME}

clean:
	go clean

test:
	go test -count 1 ./...

test-cov:
	go test -failfast -count 1 -coverpkg=./... -coverprofile=${TMPDIR}cov_profile.out ./...
	# `go env GOPATH`/bin/gcov2lcov -infile ${TMPDIR}cov_profile.out -outfile ./cov_profile.lcov

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

dep:
	go mod download
	go mod tidy

lint:
	`go env GOPATH`/bin/golangci-lint run

container-start:
	docker compose --file ./deployments/compose.yml up --detach

container-start-prod:
	docker compose --file ./deployments/compose.production.yml up --detach

container-rebuild:
	docker compose --file ./deployments/compose.yml up --detach --build

container-rebuild-prod:
	docker compose --file ./deployments/compose.production.yml up --detach --build

container-restart:
	docker compose --file ./deployments/compose.yml restart

container-restart-prod:
	docker compose --file ./deployments/compose.production.yml restart

container-stop:
	docker compose --file ./deployments/compose.yml stop

container-stop-prod:
	docker compose --file ./deployments/compose.production.yml stop

container-destroy:
	docker compose --file ./deployments/compose.yml down

container-destroy-prod:
	docker compose --file ./deployments/compose.production.yml down

container-update:
	docker compose --file ./deployments/compose.yml pull

container-update-prod:
	docker compose --file ./deployments/compose.production.yml pull

container-dev:
	docker compose --file ./deployments/compose.yml watch

container-ps:
	docker compose --file ./deployments/compose.yml ps --all

container-ps-prod:
	docker compose --file ./deployments/compose.production.yml ps --all

container-logs-all:
	docker compose --file ./deployments/compose.yml logs

container-logs-all-prod:
	docker compose --file ./deployments/compose.production.yml logs

container-logs:
	docker compose --file ./deployments/compose.yml logs go-service

container-logs-prod:
	docker compose --file ./deployments/compose.production.yml logs go-service

container-cli:
	docker compose --file ./deployments/compose.yml exec go-service bash

container-cli-prod:
	docker compose --file ./deployments/compose.production.yml exec go-service bash

container-push:
ifdef VERSION
	docker build --platform=linux/amd64 -t acikyazilim.registry.cpln.io/golang-service-template:v$(VERSION) .
	docker push acikyazilim.registry.cpln.io/golang-service-template:v$(VERSION)
else
	@echo "VERSION is not set"
endif
