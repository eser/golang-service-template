package di

import (
	"reflect"
)

// Implements is a marker type to associate implementations with interfaces.
type Implements[I any] struct{}

var Default = NewContainer() //nolint:gochecknoglobals

// Register registers an implementation instance for a given interface.
func Register(c *Container, v any) {
	c.SetValue(v)
}

func RegisterFor[I any](c *Container, v I) {
	targetType := reflect.TypeFor[I]()

	c.SetValueFor(targetType, v)
}

func RegisterFn(c *Container, fns ...any) error {
	return c.SetValuesFromFunc(fns...)
}

func Seal(c *Container) {
	c.Seal()
}

// Get retrieves a registered implementation for the given interface.
func Get[I any](c *Container) (I, bool) { //nolint:ireturn
	interfaceType := reflect.TypeFor[I]()

	impl, ok := c.Resolve(interfaceType)
	if !ok {
		var zero I

		return zero, false
	}

	return impl.Value.(I), true //nolint:forcetypeassert
}

// MustGet retrieves a registered implementation or panics if not found.
func MustGet[I any](c *Container) I { //nolint:ireturn
	interfaceType := reflect.TypeFor[I]()

	impl := c.MustResolve(interfaceType)

	return impl.Value.(I) //nolint:forcetypeassert
}

func CreateLister[I any](c *Container) func() []DependencyTarget {
	return c.CreateLister(reflect.TypeFor[I]())
}

func DynamicList[I any](c *Container) []DependencyTarget {
	interfaceType := reflect.TypeFor[I]()

	return c.DynamicList(interfaceType)
}

func CreateInvoker(c *Container, fn any) func() error {
	return c.CreateInvoker(fn)
}

func DynamicInvoke(c *Container, fn any) error {
	return c.DynamicInvoke(fn)
}
