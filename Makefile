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

.PHONY: dep-tools
dep-tools: dep # Installs tools.
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: init
init: dep-tools # Initializes the project.
	command -v pre-commit >/dev/null || brew install pre-commit
	command -v make >/dev/null || brew install make
	command -v protoc >/dev/null || brew install protobuf
	[ -f .git/hooks/pre-commit ] || pre-commit install
	command -v air >/dev/null || go install github.com/air-verse/air@latest
	command -v betteralign >/dev/null || go install github.com/dkorunic/betteralign/cmd/betteralign@latest
	command -v gcov2lcov >/dev/null || go install github.com/jandelgado/gcov2lcov@latest
	command -v goose >/dev/null || go install github.com/pressly/goose/v3/cmd/goose@latest
	command -v govulncheck >/dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	command -v mockgen >/dev/null || go install go.uber.org/mock/mockgen@latest
	command -v protoc-gen-go >/dev/null || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	command -v protoc-gen-go-grpc >/dev/null || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	command -v stringer >/dev/null || go install golang.org/x/tools/cmd/stringer@latest
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

.PHONY: sample-dev
sample-dev: ## Runs the sample service in development mode.
	air --build.bin "./tmp/samplesvc" --build.cmd "go build -o ./tmp/samplesvc ./cmd/samplesvc/"

.PHONY: sample-run
sample-run: ## Runs the sample service.
	go run ./cmd/samplesvc/

.PHONY: test
test: ## Runs the tests.
	go test -failfast -race -count 1 ./...

.PHONY: test-cov
test-cov: ## Runs the tests with coverage.
	go test -failfast -race -count 1 -coverpkg=./... -coverprofile=${TMPDIR}cov_profile.out ./...
	# gcov2lcov -infile ${TMPDIR}cov_profile.out -outfile ./cov_profile.lcov

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
	golangci-lint run ./...

.PHONY: check
check: ## Runs static analysis tools.
	govulncheck ./...
	betteralign ./...
	go vet ./...

.PHONY: fix
fix: ## Fixes code formatting and alignment.
	betteralign -apply ./...
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
	docker compose --file ./ops/docker/compose.yml exec samplesvc bash

.PHONY: container-push
container-push: ## Pushes the container to the registry.
ifdef VERSION
	docker build --platform=linux/amd64 -t acikyazilim.registry.cpln.io/samplesvc:v$(VERSION) .
	docker push acikyazilim.registry.cpln.io/samplesvc:v$(VERSION)
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
