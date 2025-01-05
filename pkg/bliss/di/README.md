# bliss/di

## Overview

The **di** package provides a lightweight yet powerful dependency injection
container for Go applications. It supports registration of concrete
implementations, interface bindings, and dependency resolution through
reflection. The package is designed to be simple to use while providing robust
dependency management capabilities.

## Features

- Type-safe dependency injection
- Interface-to-implementation binding
- Container sealing to prevent runtime modifications
- Dependency provider functions
- Dependency listing
- Support for dependency invocation

## Basic Usage

```go
import (
  ...
  "github.com/eser/go-service/pkg/bliss/di"
  ...
)

// Use the default container
container := di.Default

// Register dependencies
err := di.RegisterFn(
  container,
  configfx.RegisterDependencies,
  logfx.RegisterDependencies,
  httpfx.RegisterDependencies,
  // ... other dependencies
)
if err != nil {
  panic(err)
}

// Create and run the application
run := di.CreateInvoker(
  container,
  func(
    httpService httpfx.HttpService,
    // ... other dependencies
  ) error {
    err := httpService.Start()
    if err != nil {
      return err
    }

    return nil
  },
)

// Seal the container to prevent further modifications
di.Seal(container)

err = run()
```

## API Reference

### Registration

```go
// Register a concrete implementation
di.Register(container, &MyService{})

// Register an implementation for an interface
di.RegisterFor[MyInterface](container, &MyImplementation{})

// Register using a provider function
di.RegisterFn(container, func() (MyInterface, error) {
    return &MyImplementation{}, nil
})
```

### Resolution

```go
// Get a dependency (with error checking)
service, ok := di.Get[MyService](container)
if !ok {
    // Handle missing dependency
}

// Get a dependency (panic if not found)
service := di.MustGet[MyService](container)
```

### Dependency Listing

```go
// Create a lister function for an interface
lister := di.CreateLister[MyInterface](container)
implementations := lister()

// Get implementations directly
implementations := di.DynamicList[MyInterface](container)
```

### Invocation (preferred, w/ sealed container)

```go
// Create an invoker function
invoker := di.CreateInvoker(container, func(
    service MyService,
    repo Repository,
) error {
    // Use dependencies
    return nil
})

// Execute the invoker
err := invoker()
```

### Dynamic Invocation

```go
err := di.DynamicInvoke(container, func(
    service MyService,
    repo Repository,
) error {
    // Use dependencies
    return nil
})
```

## Example with Interfaces

```go
// Define an interface
type Adder interface {
    Add(ctx context.Context, x, y int) (int, error)
}

// Create an implementation
type AdderImpl struct {
    di.Implements[Adder] // Mark as implementation
}

func (a AdderImpl) Add(ctx context.Context, x, y int) (int, error) {
    return x + y, nil
}

// Register and use
func main() {
    container := di.NewContainer()

    // Register implementation
    di.RegisterFor[Adder](container, AdderImpl{})

    // Use through dependency injection
    err := di.DynamicInvoke(container, func(adder Adder) error {
        result, err := adder.Add(context.Background(), 2, 3)
        if err != nil {
            return err
        }
        fmt.Printf("Result: %d\n", result)
        return nil
    })
}
```

## Container Lifecycle

1. Create or use the default container
2. Register all dependencies
3. Seal the container to prevent modifications
4. Resolve and use dependencies

## Best Practices

- Register dependencies at application startup
- Seal the container before using it
- Use interfaces for better testability
- Prefer `RegisterFn` for complex initialization
- Use `di.Implements[T]` to mark interface implementations
- Handle errors from `RegisterFn` and dependency resolution

## Thread Safety

The container is not thread-safe during registration. Ensure all registrations
are complete and the container is sealed before using it concurrently.
