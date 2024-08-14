# Golang Service Template

not yet.

## Local development

Clone this repository into your projects folder

### Installation

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
  - Install golangci-lint

  ```bash
  $ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
  golangci/golangci-lint info checking GitHub for latest tag
  golangci/golangci-lint info found version: 0.0.0 for v0.0.0/os/arch64
  golangci/golangci-lint info installed ~/go/0.0.0/packages/bin/golangci-lint
  ```

  - Install gcov2lcov

  ```bash
  $ go install github.com/jandelgado/gcov2lcov@latest
  ```

  - Install pre-commit

  ```bash
  $ brew install pre-commit
  ==> Fetching dependencies for pre-commit
  ==> Fetching pre-commit
  ==> Installing dependencies for pre-commit
  ==> Installing pre-commit
  ```

- 3️⃣ Ensure that you can access private dependencies

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
  $ go mod download
  ```

- 5️⃣ Install pre-commit hooks

  ```bash
  $ pre-commit install
  ```

### Execution

#### Running the project in test mode

Before running any command, please make sure that you have configured your environment regarding your own settings first. You may found the related entries that can be configured in `.env` file.

```bash
$ go run ./cmd/service-cli/
[debug] Listening and serving HTTP on :8080
```

- You can access http://localhost:8080/ to check if the project is running

### Testing the project

```bash
$ go test -v --coverpkg=./... ./...
```

### Development (with Docker Compose)

```bash
# first start freshly built containers in background (daemon mode)
$ docker compose up --build -d

# then launch watch mode to reflect changes on time
$ docker compose watch
```

### Contribution

- Create a new branch for your feature
- Make your changes
- Test your changes (see [Execution](#execution) below)
- Commit your changes
- Push your changes to your branch
- Create a pull request to the `main` branch
- Wait for review and merge
- Delete your branch


## More Information

This project is bootstrapped from https://github.com/eser/golang-service-template.
See the source repository for further details.
