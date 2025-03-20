When working with this codebase, please follow these architectural principles and guidelines:

1. Hexagonal Architecture
- The codebase follows strict hexagonal architecture (ports and adapters) principles
- Business logic in pkg/sample/business/ must not have external dependencies
- All external interactions must go through interfaces (ports) defined in the business layer
- Adapters in pkg/sample/adapters/ implement these interfaces
- Ports exist as interfaces in the business layer

1. Dependency Rules
- Business logic can only depend on other business logic
- Adapters can only depend on business interfaces they implement
- No circular dependencies allowed
- Entry points (cmd/) wire everything together using dependency injection

1. Package Structure
- Business logic packages should define their interfaces and types first
- Each business domain should have its own package under pkg/sample/business/
- Adapters should be organized by technology (http, storage, etc.)
- Configuration should be injected via the appcontext

1. Conventions
- Use snake_case for file names.
- Follow idiomatic Go conventions.

1. Code Generation
- Database code is generated using sqlc
- Maintain clean separation between generated and hand-written code

1. Error Handling
- Business errors should be defined in the business layer
- Wrap external errors before returning them
- Use meaningful error types and messages
- Avoid panic as much as possible (especially in business logic)

1. Testing
- Business logic must have unit tests
- Adapters should have integration tests
- Use interfaces for mocking in tests
- Test files should be next to their tested code with package names prefixed with _test

1. Configuration
- Use environment variables for runtime configuration
- Keep secrets out of the codebase
- Use .env.local and config.local.json for development overrides

1. Database
- All SQL queries should be in etc/data/{datasource_name}/queries/
- Use migrations for schema changes
- Keep migrations forward-only
- Document schema changes

Remember:
- Business logic is sacred - keep it clean and dependency-free
- All external interactions must go through well-defined interfaces
- Maintain clear separation between layers
