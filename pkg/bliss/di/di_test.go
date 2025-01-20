package di_test

import (
	"context"
	"testing"

	"github.com/eser/go-service/pkg/bliss/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define test interfaces and implementations.
type Adder interface {
	Add(ctx context.Context, x int, y int) (int, error)
}

type AdderImpl struct {
	di.Implements[Adder]
}

func (AdderImpl) Add(_ context.Context, x int, y int) (int, error) {
	return x + y, nil
}

type Multiplier interface {
	Multiply(ctx context.Context, x int, y int) (int, error)
}

type MultiplierImpl struct {
	di.Implements[Multiplier]
}

func (MultiplierImpl) Multiply(_ context.Context, x int, y int) (int, error) {
	return x * y, nil
}

func TestRegisterAndGetValue(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register the implementation
	di.Register(c, AdderImpl{}) //nolint:exhaustruct
	// di.RegisterFor[Adder](c, AdderImpl{}) //nolint:exhaustruct

	// Get the implementation
	adder, ok := di.Get[Adder](c)
	require.True(t, ok, "Expected to retrieve Adder implementation")
	require.NotNil(t, adder, "Adder implementation should not be nil")

	// Use the implementation
	result, err := adder.Add(context.Background(), 2, 3)
	require.NoError(t, err, "Adder.Add should not return an error")
	assert.Equal(t, 5, result, "Adder.Add should return the correct sum")
}

func TestRegisterProviderAndGetValue(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register the implementation using provider
	err := di.RegisterFn(c, func() (Adder, error) {
		return AdderImpl{}, nil //nolint:exhaustruct
	})
	require.NoError(t, err, "RegisterFn should not return an error")

	// Get the implementation
	adder, ok := di.Get[Adder](c)
	require.True(t, ok, "Expected to retrieve Adder implementation")
	require.NotNil(t, adder, "Adder implementation should not be nil")

	// Use the implementation
	result, err := adder.Add(context.Background(), 2, 3)
	require.NoError(t, err, "Adder.Add should not return an error")
	assert.Equal(t, 5, result, "Adder.Add should return the correct sum")
}

func TestMustGetValue(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register the implementation
	di.RegisterFor[Multiplier](c, MultiplierImpl{}) //nolint:exhaustruct

	// MustGet should return the implementation
	multiplier := di.MustGet[Multiplier](c)
	require.NotNil(t, multiplier, "Multiplier implementation should not be nil")

	// Use the implementation
	result, err := multiplier.Multiply(context.Background(), 2, 3)
	require.NoError(t, err, "Multiplier.Multiply should not return an error")
	assert.Equal(t, 6, result, "Multiplier.Multiply should return the correct product")
}

func TestMustGetValueWithProvider(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register the implementation using provider
	err := di.RegisterFn(c, func() (Multiplier, error) {
		return MultiplierImpl{}, nil //nolint:exhaustruct
	})
	require.NoError(t, err, "RegisterFn should not return an error")

	// MustGet should return the implementation
	multiplier := di.MustGet[Multiplier](c)
	require.NotNil(t, multiplier, "Multiplier implementation should not be nil")

	// Use the implementation
	result, err := multiplier.Multiply(context.Background(), 2, 3)
	require.NoError(t, err, "Multiplier.Multiply should not return an error")
	assert.Equal(t, 6, result, "Multiplier.Multiply should return the correct product")
}

func TestMustGetValue_Panic(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Ensure that MustGet panics if not registered
	assert.PanicsWithValue(t,
		"No implementation registered for type di_test.Adder",
		func() {
			di.MustGet[Adder](c)
		},
		"Expected MustGet to panic when the implementation is not registered",
	)
}

func TestMustGetValue_PanicWithProvider(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register a failing provider
	err := di.RegisterFn(c, func() (Adder, error) {
		return nil, assert.AnError
	})
	assert.ErrorIs(t, err, assert.AnError, "RegisterFn should return the provider error")

	// Ensure that MustGet panics if provider returns error
	assert.PanicsWithValue(t,
		"No implementation registered for type di_test.Adder",
		func() {
			di.MustGet[Adder](c)
		},
		"Expected MustGet to panic when the provider returns an error",
	)
}

func TestGetValue_NotRegistered(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Ensure that Get returns false if not registered
	adder, ok := di.Get[Adder](c)
	assert.False(t, ok, "Expected Get to return false when not registered")
	assert.Zero(t, adder, "Adder should be zero value when not registered")
}

func TestGetValue_NotRegisteredWithProvider(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register a failing provider
	err := di.RegisterFn(c, func() (Adder, error) {
		return nil, assert.AnError
	})
	assert.ErrorIs(t, err, assert.AnError, "RegisterFn should return the provider error")

	// Ensure that Get returns false if provider returns error
	adder, ok := di.Get[Adder](c)
	assert.False(t, ok, "Expected Get to return false when provider returns error")
	assert.Zero(t, adder, "Adder should be zero value when provider returns error")
}

func TestInvoke(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register implementations
	di.RegisterFor[Adder](c, AdderImpl{})           //nolint:exhaustruct
	di.RegisterFor[Multiplier](c, MultiplierImpl{}) //nolint:exhaustruct

	// Invoke a function that takes dependencies
	err := di.DynamicInvoke(c, func(adder Adder, multiplier Multiplier) error {
		sum, err := adder.Add(context.Background(), 2, 3)
		require.NoError(t, err, "Adder.Add should not return an error")
		assert.Equal(t, 5, sum, "Adder.Add should return the correct sum")

		product, err := multiplier.Multiply(context.Background(), 2, 3)
		require.NoError(t, err, "Multiplier.Multiply should not return an error")
		assert.Equal(t, 6, product, "Multiplier.Multiply should return the correct product")

		return nil
	})
	require.NoError(t, err, "DynamicInvoke should not return an error")
}

func TestInvoke_MissingDependency(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register only the Adder implementation
	di.RegisterFor[Adder](c, AdderImpl{}) //nolint:exhaustruct

	// Ensure that Invoke panics if a dependency is missing
	assert.PanicsWithValue(t,
		"No implementation registered for type di_test.Multiplier",
		func() {
			di.DynamicInvoke(c, func(multiplier Multiplier) error {
				// This function should not be called
				return nil
			})
		},
		"Expected Invoke to panic when a dependency is not registered",
	)
}

func TestInvoke_NonFunction(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Ensure that Invoke panics if passed a non-function
	assert.PanicsWithValue(t,
		"DynamicInvoke parameter must be a function",
		func() {
			di.DynamicInvoke(c, 42) // Passing a non-function
		},
		"Expected Invoke to panic when a non-function is passed",
	)
}

func TestSeal(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register the implementation
	di.RegisterFor[Adder](c, AdderImpl{}) //nolint:exhaustruct

	// Seal the container
	di.Seal(c)

	// Ensure that registering after sealing panics
	assert.Panics(t, func() {
		di.Register(c, MultiplierImpl{}) //nolint:exhaustruct
	}, "Expected Register to panic after container is sealed")

	// Verify that existing registrations still work
	adder := di.MustGet[Adder](c)
	require.NotNil(t, adder, "Should still be able to get registered implementations after sealing")
}

func TestCreateLister(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register implementations
	di.RegisterFor[Adder](c, AdderImpl{})           //nolint:exhaustruct
	di.RegisterFor[Multiplier](c, MultiplierImpl{}) //nolint:exhaustruct

	// Create lister for Adder
	lister := di.CreateLister[Adder](c)

	// Ensure that lister returns the correct implementations
	implementations := lister()
	require.Len(t, implementations, 1, "Expected one Adder implementation")
	assert.IsType(t, AdderImpl{}, implementations[0].Value, "Expected AdderImpl type")

	// Verify the implementation works
	adder := implementations[0].Value.(Adder)
	result, err := adder.Add(context.Background(), 2, 3)
	require.NoError(t, err)
	assert.Equal(t, 5, result)
}

func TestDynamicList(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register implementations
	di.RegisterFor[Adder](c, AdderImpl{})           //nolint:exhaustruct
	di.RegisterFor[Multiplier](c, MultiplierImpl{}) //nolint:exhaustruct

	// Get dynamic list for Adder
	implementations := di.DynamicList[Adder](c)

	// Ensure that dynamic list returns the correct implementations
	require.Len(t, implementations, 1, "Expected one Adder implementation")
	assert.IsType(t, AdderImpl{}, implementations[0].Value, "Expected AdderImpl type")

	// Verify the implementation works
	adder := implementations[0].Value.(Adder)
	result, err := adder.Add(context.Background(), 2, 3)
	require.NoError(t, err)
	assert.Equal(t, 5, result)
}

func TestCreateInvoker(t *testing.T) {
	t.Parallel()

	c := di.NewContainer()

	// Register implementations
	di.RegisterFor[Adder](c, AdderImpl{})           //nolint:exhaustruct
	di.RegisterFor[Multiplier](c, MultiplierImpl{}) //nolint:exhaustruct

	// Create invoker for a function that takes dependencies
	invoker := di.CreateInvoker(c, func(adder Adder, multiplier Multiplier) error {
		sum, err := adder.Add(context.Background(), 2, 3)
		require.NoError(t, err, "Adder.Add should not return an error")
		assert.Equal(t, 5, sum, "Adder.Add should return the correct sum")

		product, err := multiplier.Multiply(context.Background(), 2, 3)
		require.NoError(t, err, "Multiplier.Multiply should not return an error")
		assert.Equal(t, 6, product, "Multiplier.Multiply should return the correct product")

		return nil
	})

	// Ensure that invoker does not return an error
	require.NoError(t, invoker(), "Invoker should not return an error")
}
