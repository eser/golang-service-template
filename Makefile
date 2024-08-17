.PHONY: dev build multiarch-build run run-docker clean test test-cov test-ci dep lint container-start container-rebuild container-restart container-stop container-destroy container-update container-dev container-ps container-logs-all container-logs container-cli
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

container-rebuild:
	docker compose --file ./deployments/compose.yml up --detach --build

container-restart:
	docker compose --file ./deployments/compose.yml restart

container-stop:
	docker compose --file ./deployments/compose.yml stop

container-destroy:
	docker compose --file ./deployments/compose.yml down

container-update:
	docker compose --file ./deployments/compose.yml pull

container-dev:
	docker compose --file ./deployments/compose.yml watch

container-ps:
	docker compose --file ./deployments/compose.yml ps --all

container-logs-all:
	docker compose --file ./deployments/compose.yml logs

container-logs:
	docker compose --file ./deployments/compose.yml logs go-service

container-cli:
	docker compose --file ./deployments/compose.yml exec go-service bash
