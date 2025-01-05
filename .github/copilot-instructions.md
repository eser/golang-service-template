This codebase consists of several interconnected Golang and JavaScript projects that share a common root folder. All packages and application modules, regardless of language or runtime, are located under the `pkg/` directory. The `cmd/` directory includes codebase tooling and console commands. The `ops/` directory contains specifications and configurations for Kubernetes, Docker, and other tools.

For more information, refer to `ARCHITECTURE.md` located in the `.github/` directory. Before adding code to this codebase using generative AI, you must thoroughly understand the existing structure. Changes should be made without disrupting design decisions and must maintain consistency. Preserving conventions is of utmost importance.

You should obtain URLs, namespaces, and other relevant information about existing modules and packages from the `go.mod` and `sqlc.yaml` files. Additionally, environment variables and application configurations in Go can be accessed through the `configfx` module, and database connections obtained via configuration can be used through the `datafx` packages. Please use the same namespaces, conventions, and URLs consistently when adding or refactoring modules or packages.

When generating or creating new components, you must refer to definitions like `samplesvc` and ensure that the necessary additions are made throughout the entire codebase, covering all relevant layers (full-stack). There should be no missing definitions, and code generation must always be comprehensive. Without additional input, always continue generating additional files such as service implementations, repositories, routes/procedures, entry points, etc. In the code sections, instead of placeholders like "to be added later" or "this will come here," the code should be as close as possible to production-ready, clearly utilizing existing objects. Even the simplest methods should be fully implemented; nothing should be left as a placeholder.

### Service Modules

Each **service module** is defined under the `pkg/` directory (e.g., `sample` is our sample service).

Each service module should be suffixed with "svc" or "http" depending on its protocol and should have its own `sqlc` definitions, repositories, services, and other components. Service modules have their own entry points defined in the `cmd/` directory with the same name and corresponding `Makefile` targets for running the service locally. After creating a service module, you must update the unit tests and the README.

### Domain Entities

Each **domain entity** (or **domain object**) belongs to a business layer. Because domain entities are exposed to different services, each should define RPC methods for gRPC, route definitions for HTTP, and the necessary Data Definition Language (DDL) and Data Manipulation Language (DML) definitions in `sqlc` for the database. The `sqlc` configuration can be found in `sqlc.yaml`, and existing configurations must be considered. Each `sqlc` definition is organized per service module; if a service module already has its own `sqlc` definitions, you should use them; otherwise, create a new `sqlc` definition for that service module. Additionally, create repositories that use these generated `sqlc` methods, along with services that utilize unit-of-work patterns with these repositories.

For any new domain entities created, you must update the documentation in the `docs/` folder, the unit tests and the README.

### Conventions

**Language Independent**:
- Use kebab-case for file names.

**Golang**:
- Follow idiomatic Go conventions.
