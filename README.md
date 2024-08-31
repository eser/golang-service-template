# Golang Service Template

`Golang Service Template` project is designed to provide a robust foundation that is always ready to be open-sourced, accelerating development and fostering a unified understanding across disciplines. It empowers teams to quickly adopt best practices and streamline the project setup, ensuring consistency and clarity from the very start.


## Local development

Clone this repository into your projects folder


### Structure

This project inherits the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) structure but includes its own interpretation.


### Directories

- The `pkg` directory contains packages for project modules. The absence of a direct `internal` folder helps us maintain a mentally isolated modular structure.
- The `deployment` directory contains local and remote infrastructure configuration files, such as `compose.yml`.
- The `cmd` directory contains the entrypoint for the project binaries.


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

  **On macOS and Homebrew (automatically):**

  ```bash
  $ make init
  ```

  If it fails on any step, you can install them manually by following the steps below. Otherwise, you can skip the rest of the steps.


  **On other OS or without Homebrew:**

  - Install and enable [pre-commit](https://pre-commit.com/#install)
  - Install [GNU make](https://www.gnu.org/software/make/)
  - Install [govulncheck](https://go.googlesource.com/vuln)
  - Install [gcov2lcov](https://github.com/jandelgado/gcov2lcov?tab=readme-ov-file#installation)

  ```bash
  $ brew install pre-commit
  ==> Fetching dependencies for pre-commit
  ==> Fetching pre-commit
  ==> Installing dependencies for pre-commit
  ==> Installing pre-commit
  ...

  $ brew install make
  ==> Fetching dependencies for make
  ==> Fetching make
  ==> Installing dependencies for make
  ==> Installing make
  ...

  $ pre-commit install
  pre-commit installed at .git/hooks/pre-commit

  $ go install golang.org/x/vuln/cmd/govulncheck@latest
  go: downloading golang.org/x/vuln/cmd/govulncheck v0.0.0

  $ go install github.com/jandelgado/gcov2lcov@latest
  go: downloading github.com/jandelgado/gcov2lcov v0.0.0
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

### Execution

#### Running the project

Before running any command, please make sure that you have configured your environment regarding your own settings first. You may found the related entries that can be configured in `.env` file.

```bash
$ make run
17:05:20.311 INFO HttpService is starting... {"addr":":8080"}
```

- You can access http://localhost:8080/ to check if the project is running


### Running the project (with hot-reloading development mode)

```bash
$ make dev
```


### Testing the project

```bash
$ make test
```


### Development (with Docker Compose)

```bash
# first start freshly built containers in background (daemon mode)
$ make container-start

# then launch watch mode to reflect changes on time
$ make container-dev
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
