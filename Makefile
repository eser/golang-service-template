# .RECIPEPREFIX := $(.RECIPEPREFIX)<space>
TESTCOVERAGE_THRESHOLD=0

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

.PHONY: init
init:
	command -v deno >/dev/null || curl -fsSL https://deno.land/install.sh | sh
	command -v pre-commit >/dev/null || brew install pre-commit
	command -v make >/dev/null || brew install make
	command -v protoc >/dev/null || brew install protobuf
	[ -f .git/hooks/pre-commit ] || pre-commit install
	command -v air >/dev/null || go install github.com/air-verse/air@latest
	command -v govulncheck >/dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	command -v betteralign >/dev/null || go install github.com/dkorunic/betteralign/cmd/betteralign@latest
	command -v goose >/dev/null || go install github.com/pressly/goose/v3/cmd/goose@latest
	command -v gcov2lcov >/dev/null || go install github.com/jandelgado/gcov2lcov@latest
	command -v protoc-gen-go >/dev/null || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	command -v protoc-gen-go-grpc >/dev/null || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: dev-identitysvc
dev-identitysvc:
	air --build.bin "./tmp/identitysvc" --build.cmd "go build -o ./tmp/identitysvc ./cmd/identitysvc/"

.PHONY: dev-broadcastsvc
dev-broadcastsvc:
	air --build.bin "./tmp/broadcastsvc" --build.cmd "go build -o ./tmp/broadcastsvc ./cmd/broadcastsvc/"

.PHONY: dev-broadcasthttp
dev-broadcasthttp:
	air --build.bin "./tmp/broadcasthttp" --build.cmd "go build -o ./tmp/broadcasthttp ./cmd/broadcasthttp/"

.PHONY: dev-testhttp
dev-testhttp:
	air --build.bin "./tmp/testhttp" --build.cmd "go build -o ./tmp/testhttp ./cmd/testhttp/"

.PHONY: run-identitysvc
run-identitysvc:
	go run ./cmd/identitysvc/

.PHONY: run-broadcastsvc
run-broadcastsvc:
	go run ./cmd/broadcastsvc/

.PHONY: run-broadcasthttp
run-broadcasthttp:
	go run ./cmd/broadcasthttp/

.PHONY: run-testhttp
run-testhttp:
	go run ./cmd/testhttp/

.PHONY: migrate
migrate:
	go run ./cmd/migrate/ $(ARGS)

.PHONY: build
build:
	go build -v ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	go clean

.PHONY: test-api
test-api:
	cd ./ops/api-tests/ && \
	kreyac environment set-active development --disable-telemetry && \
	kreyac collection invoke '**' --disable-telemetry && \
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

.PHONY: check
check:
	`go env GOPATH`/bin/govulncheck ./...
	`go env GOPATH`/bin/betteralign ./...

.PHONY: postgres-start
postgres-start:
	docker compose --file ./ops/docker/compose.yml up --detach postgres

.PHONY: postgres-stop
postgres-stop:
	docker compose --file ./ops/docker/compose.yml stop postgres

.PHONY: container-start
container-start:
	docker compose --file ./ops/docker/compose.yml up --detach

.PHONY: container-rebuild
container-rebuild:
	docker compose --file ./ops/docker/compose.yml up --detach --build

.PHONY: container-restart
container-restart:
	docker compose --file ./ops/docker/compose.yml restart

.PHONY: container-stop
container-stop:
	docker compose --file ./ops/docker/compose.yml stop

.PHONY: container-destroy
container-destroy:
	docker compose --file ./ops/docker/compose.yml down

.PHONY: container-update
container-update:
	docker compose --file ./ops/docker/compose.yml pull

.PHONY: container-dev
container-dev:
	docker compose --file ./ops/docker/compose.yml watch

.PHONY: container-ps
container-ps:
	docker compose --file ./ops/docker/compose.yml ps --all

.PHONY: container-logs
container-logs:
	docker compose --file ./ops/docker/compose.yml logs

.PHONY: container-cli
container-cli:
	docker compose --file ./ops/docker/compose.yml exec identitysvc bash

.PHONY: container-push
container-push:
ifdef VERSION
	docker build --platform=linux/amd64 -t acikyazilim.registry.cpln.io/identitysvc:v$(VERSION) .
	docker push acikyazilim.registry.cpln.io/identitysvc:v$(VERSION)
else
	@echo "VERSION is not set"
endif

.PHONY: generate-proto
# --ts_proto_opt="context=true,env=node,lowerCaseServiceMethods=true,outputServices=grpc-js,removeEnumPrefix=true,snakeToCamel=true,useAbortSignal=true,useAsyncIterable=true,useReadonlyTypes=true,comments=false,useNullAsOptional=true"
generate-proto:
	@{ \
	  for f in ./specs/proto/*; do \
	    current_proto="$$(basename $$f)"; \
	    echo "Generating stubs for $$current_proto"; \
			\
			protoc --plugin=./pkg/web/node_modules/.bin/protoc-gen-ts_proto \
				--proto_path=./specs/proto/ \
				--ts_proto_out=./pkg/web/proto-ts/ \
				--ts_proto_opt="context=true,lowerCaseServiceMethods=true,outputServices=grpc-js,removeEnumPrefix=false,snakeToCamel=true,useReadonlyTypes=true,comments=false,useNullAsOptional=true" \
				--go_out=./pkg/proto-go/ --go_opt=paths=source_relative \
				--go-grpc_out=./pkg/proto-go/ --go-grpc_opt=paths=source_relative \
				"./specs/proto/$$current_proto/$$current_proto.proto"; \
	  done \
	}

%:
	@:
