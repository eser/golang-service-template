# Golang Service Template

`Golang Service Template` project is designed to provide a robust foundation
that is always ready to be open-sourced, accelerating development and fostering
a unified understanding across disciplines. It empowers teams to quickly adopt
best practices and streamline the project setup, ensuring consistency and
clarity from the very start.


## Structure

This project inherits the
[Standard Go Project Layout](https://github.com/golang-standards/project-layout)
structure but includes its own interpretation.

The template code also serves as a service template, implemented with Hexagonal
Architecture to support typical software organizations striving to apply
Event-Driven Architecture (EDA) and Domain-Driven Design (DDD).

The decision to use Hexagonal Architecture is based on its simplicity as one of
the most straightforward implementations of Onion Architecture, with which I
have extensive experience. Additionally, it is flexible enough to evolve into
more structured or complex systems, such as Clean Architecture, as project
requirements grow.


## Hexagonal Architecture Overview

This project follows hexagonal architecture principles, also known as ports and adapters pattern:

### Business Logic (`pkg/sample/business/`)
- Contains domain entities and business rules
- Defines interfaces (ports) that the outside world must implement
- No external dependencies, pure business logic
- Example: `channels.Service` interface in `pkg/sample/business/channels/service.go`

### Adapters (`pkg/sample/adapters/`)
- Implements interfaces defined by the business logic
- Handles external concerns (HTTP, GRPC, database, etc.)
- Organized by technology/concern:
  - `appcontext/`: Application configuration and context
  - `storage/`: Database repositories
  - `http/`: HTTP server and handlers
  - `grpc/`: GRPC server and handlers


## Directories

```
.
├── cmd/                     # Application entry points
│   ├── migrate/             # Database migration tool (based on goose)
│   ├── manage/              # Manages codebase and app-related resources
│   └── serve/               # Main service entry point
├── pkg/
│   └── sample/              # Our application code
│       ├── adapters/        # Implementation of ports (adapters)
│       │   ├── appcontext/  # Application context and configuration
│       │   ├── http/        # HTTP server and handlers
│       │   ├── grpc/        # GRPC server and handlers
│       │   └── storage/     # Database repositories
│       └── business/        # Business logic and domain models
│           ├── channels/    # Channel-related business objects
│           └── tenants/     # Tenant-related business objects
├── etc/
│   └── data/                # Data source-related files
│       └── default/         # Data source name
│           ├── migrations/  # SQL migration files
│           └── queries/     # SQL query definitions
└── ops/                     # Operational configurations
    └── docker/              # Docker-related files
```

## Installation

- 1️⃣ Ensure that `golang` tools are _properly_ installed on your machine

  ```bash
  $ go version
  go version go0.0.0 os/arch64

  $ go env GOPATH
  ~/go/0.0.0/packages
  ```

  **Important rules:**
  - `GOPATH` envvar must be set to output of `go env GOPATH`
  - `$GOPATH/bin` must be included to your `PATH` envvar

  So, your `~/.zprofile` should include something like this:

  ```
  export GOPATH="$(go env GOPATH)"
  export PATH="$PATH:$GOPATH/bin"
  ```

- 2️⃣ Install prerequisites

  **On macOS and Homebrew (automatically):**

  ```bash
  $ make init
  ```

  If it fails on any step, you can install them manually by following the steps
  below. Otherwise, you can skip the rest of the steps.

  **On other OS or without Homebrew:**

  - Install and enable [pre-commit](https://pre-commit.com/#install)
  - Install [GNU make](https://www.gnu.org/software/make/)
  - Install [act](https://nektosact.com/installation/index.html)
  - Install [protobuf](https://github.com/protocolbuffers/protobuf/releases)
  - Install [Air](https://github.com/air-verse/air#installation)
  - Install generators and checkers via `go tool`

  ```bash
  $ brew install pre-commit
  ==> Fetching dependencies for pre-commit
  ==> Fetching pre-commit
  ==> Installing dependencies for pre-commit
  ==> Installing pre-commit
  ...

  $ pre-commit install
  pre-commit installed at .git/hooks/pre-commit

  $ brew install make
  ==> Fetching dependencies for make
  ==> Fetching make
  ==> Installing dependencies for make
  ==> Installing make
  ...

  $ brew install act
  ==> Fetching dependencies for act
  ==> Fetching act
  ==> Installing dependencies for act
  ==> Installing act
  ...

  $ brew install protobuf
  ==> Fetching dependencies for protobuf
  ==> Fetching protobuf
  ==> Installing dependencies for protobuf
  ==> Installing protobuf
  ...

  $ make init-generators
  go: downloading github.com/sqlc-dev/sqlc/cmd/sqlc v0.0.0
  go: downloading github.com/vektra/mockery/v2 v0.0.0
  go: downloading golang.org/x/tools/cmd/stringer v0.0.0
  go: downloading github.com/jandelgado/gcov2lcov v0.0.0
  go: downloading google.golang.org/protobuf/cmd/protoc-gen-go v0.0.0
  go: downloading google.golang.org/grpc/cmd/protoc-gen-go-grpc v0.0.0

  $ make init-checkers
  go: downloading github.com/golangci/golangci-lint/cmd/golangci-lint v0.0.0
  go: downloading github.com/dkorunic/betteralign/cmd/betteralign v0.0.0
  go: downloading golang.org/x/vuln/cmd/govulncheck v0.0.0
  ```

- 3️⃣ (Optional) Ensure that you can access private dependencies

  You need to get a Personal Access Token from your GitHub account in order to
  download private dependencies.

  To get these, visit https://github.com/settings/tokens/new and create a new
  token with the `read:packages` scope.

  Then, you need to create or edit the `.netrc` file in your home directory with
  the following content:

  ```
  machine github.com login <your-github-username> password <your-github-access-token>
  ```

- 4️⃣ Download required modules for the project

  ```bash
  $ make dep
  ```

- 5️⃣ (Optional) Check installation via pre-commit scripts

  ```bash
  $ pre-commit run --all-files
  ```

## Execution

### Running the project

Before running any command, please make sure that you have configured your
environment regarding your own settings first. You may found the related entries
that can be configured in `.env` file.

```bash
$ make run
02:27:53.026 INFO adding datasource connection {"name":"default","dialect":"postgres"}
02:27:53.563 INFO successfully added datasource connection {"name":"default"}
02:27:53.563 INFO Starting service {"name":"sample","environment":"development","features":{"Dummy":true}}
02:27:53.564 INFO HttpService is starting... {"addr":":8080"}
```

- You can access http://localhost:8080/ to check if the project is running

## Running the project (with hot-reloading development mode)

```bash
$ make dev
```

## Testing the project

```bash
$ make test
```

## Development (with Docker Compose)

```bash
# first start freshly built containers in background (daemon mode)
$ make container-start

# then launch watch mode to reflect changes on time
$ make container-dev
```

## Contribution

- Create a new branch for your feature
- Make your changes
- Test your changes (see [Execution](#execution) below)
- Commit your changes
- Push your changes to your branch
- Create a pull request to the `main` branch
- Wait for review and merge
- Delete your branch


## More Information

This project is bootstrapped from
https://github.com/eser/golang-service-template. See the source repository for
further details.


## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
