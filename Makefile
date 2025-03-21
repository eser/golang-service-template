# .RECIPEPREFIX := $(.RECIPEPREFIX)<space>
TESTCOVERAGE_THRESHOLD=0

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

default: help

.PHONY: help
help: ## Shows help for each of the Makefile recipes.
	@echo 'Commands:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: dep
dep: ## Downloads dependencies.
	go mod download
	go mod tidy

.PHONY: init-tools
init-tools: ## Initializes tools.
	command -v pre-commit >/dev/null || brew install pre-commit
	[ -f .git/hooks/pre-commit ] || pre-commit install
	command -v make >/dev/null || brew install make
	command -v act >/dev/null || brew install act
	command -v protoc >/dev/null || brew install protobuf
	go tool -n air >/dev/null || go get -tool github.com/air-verse/air@latest

.PHONY: init-generators
init-generators: ## Initializes generators.
	go tool -n sqlc >/dev/null || go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go tool -n mockery > /dev/null || go get -tool github.com/vektra/mockery/v2@latest
	go tool -n stringer >/dev/null || go get -tool golang.org/x/tools/cmd/stringer@latest
	go tool -n gcov2lcov >/dev/null || go get -tool github.com/jandelgado/gcov2lcov@latest
	go tool -n protoc-gen-go >/dev/null || go get -tool google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go tool -n protoc-gen-go-grpc >/dev/null || go get -tool google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: init-checkers
init-checkers: ## Initializes checkers.
	go tool -n golangci-lint >/dev/null || go get -tool github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go tool -n betteralign >/dev/null || go get -tool github.com/dkorunic/betteralign/cmd/betteralign@latest
	go tool -n govulncheck >/dev/null || go get -tool golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: init
init: init-tools init-generators init-checkers dep # Initializes the project.
	# cp -n .env.example .env || true

.PHONY: generate
generate: ## Runs auto-generated code generation tools.
	go generate ./...

.PHONY: migrate
migrate: ## Runs the migration command.
	go run ./cmd/migrate/ $(ARGS)

.PHONY: build
build: ## Builds the entire codebase.
	go build -v ./...

.PHONY: clean
clean: ## Cleans the entire codebase.
	go clean

.PHONY: dev
dev: ## Runs the sample service in development mode.
	go tool air --build.bin "./tmp/serve" --build.cmd "go build -o ./tmp/serve ./cmd/serve/"

.PHONY: run
run: ## Runs the sample service.
	go run ./cmd/serve/

.PHONY: test
test: ## Runs the tests.
	go test -failfast -race -count 1 ./...

.PHONY: test-cov
test-cov: ## Runs the tests with coverage.
	go test -failfast -race -count 1 -coverpkg=./... -coverprofile=${TMPDIR}cov_profile.out ./...
	# go tool gcov2lcov -infile ${TMPDIR}cov_profile.out -outfile ./cov_profile.lcov

.PHONY: test-view-html
test-view-html: ## Views the test coverage in HTML.
	go tool cover -html ${TMPDIR}cov_profile.out -o ${TMPDIR}cov_profile.html
	open ${TMPDIR}cov_profile.html

.PHONY: test-ci
test-ci: test-cov # Runs the tests with coverage and check if it's above the threshold.
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

.PHONY: lint
lint: ## Runs the linting command.
	go tool golangci-lint run ./...

.PHONY: check
check: ## Runs static analysis tools.
	go tool govulncheck ./...
	go tool betteralign ./...
	go vet ./...

.PHONY: fix
fix: ## Fixes code formatting and alignment.
	go tool betteralign -apply ./...
	go fmt ./...

.PHONY: postgres-start
postgres-start: ## Starts the postgres container.
	docker compose --file ./ops/docker/compose.yml up --detach postgres

.PHONY: postgres-stop
postgres-stop: ## Stops the postgres container.
	docker compose --file ./ops/docker/compose.yml stop postgres

.PHONY: container-start
container-start: ## Starts the container.
	docker compose --file ./ops/docker/compose.yml up --detach

.PHONY: container-rebuild
container-rebuild: ## Rebuilds the container.
	docker compose --file ./ops/docker/compose.yml up --detach --build

.PHONY: container-restart
container-restart: ## Restarts the container.
	docker compose --file ./ops/docker/compose.yml restart

.PHONY: container-stop
container-stop: ## Stops the container.
	docker compose --file ./ops/docker/compose.yml stop

.PHONY: container-destroy
container-destroy: ## Destroys the container.
	docker compose --file ./ops/docker/compose.yml down

.PHONY: container-update
container-update: ## Updates the container.
	docker compose --file ./ops/docker/compose.yml pull

.PHONY: container-dev
container-dev: ## Watches the container.
	docker compose --file ./ops/docker/compose.yml watch

.PHONY: container-ps
container-ps: ## Lists all containers.
	docker compose --file ./ops/docker/compose.yml ps --all

.PHONY: container-logs
container-logs: ## Shows the logs of the container.
	docker compose --file ./ops/docker/compose.yml logs

.PHONY: container-cli
container-cli: ## Opens a shell in the container.
	docker compose --file ./ops/docker/compose.yml exec sample bash

.PHONY: container-push
container-push: ## Pushes the container to the registry.
ifdef VERSION
	docker build --platform=linux/amd64 -t acikyazilim.registry.cpln.io/sample:v$(VERSION) .
	docker push acikyazilim.registry.cpln.io/sample:v$(VERSION)
else
	@echo "VERSION is not set"
endif

.PHONY: generate-proto
generate-proto: ## Generates the proto stubs.
	@{ \
	  for f in ./specs/proto/*; do \
	    current_proto="$$(basename $$f)"; \
	    echo "Generating stubs for $$current_proto"; \
			\
			protoc --proto_path=./specs/proto/ \
				--go_out=./pkg/proto-go/ --go_opt=paths=source_relative \
				--go-grpc_out=./pkg/proto-go/ --go-grpc_opt=paths=source_relative \
				"./specs/proto/$$current_proto/$$current_proto.proto"; \
	  done \
	}

%:
	@:
