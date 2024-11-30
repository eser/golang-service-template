package di

import (
	"fmt"
	"reflect"
)

type Provider func(args []any) any

// Container interface defines the methods for dependency injection container.
type Container interface {
	SetValue(value any)
	SetValueFor(interfaceType reflect.Type, value any)
	SetValuesFromFunc(fns ...any) error

	Seal()

	Resolve(interfaceType reflect.Type) (DependencyTarget, bool)
	MustResolve(interfaceType reflect.Type) DependencyTarget

	CreateLister(interfaceType reflect.Type) func() []DependencyTarget
	DynamicList(interfaceType reflect.Type) []DependencyTarget

	CreateInvoker(fn any) func() error
	DynamicInvoke(fn any) error
}

type DependencyTarget struct {
	Value           any
	ReflectionValue reflect.Value
}

// ContainerImpl is the concrete implementation of the Container interface.
type ContainerImpl struct {
	dependencies map[reflect.Type]DependencyTarget

	isSealed bool
}

var _ Container = (*ContainerImpl)(nil)

var (
	reflectTypeError     = reflect.TypeFor[error]()     //nolint:gochecknoglobals
	reflectTypeContainer = reflect.TypeFor[Container]() //nolint:gochecknoglobals
)

// NewContainer creates a new dependency injection container.
func NewContainer() *ContainerImpl {
	return &ContainerImpl{
		isSealed: false,

		dependencies: make(map[reflect.Type]DependencyTarget),
	}
}

func (c *ContainerImpl) SetValue(value any) {
	if c.isSealed {
		panic("Container is sealed")
	}

	reflectionValue := reflect.ValueOf(value)
	reflectionType := reflectionValue.Type()

	c.dependencies[reflectionType] = DependencyTarget{
		ReflectionValue: reflectionValue,
		Value:           value,
	}
}

func (c *ContainerImpl) SetValueFor(interfaceType reflect.Type, value any) {
	if c.isSealed {
		panic("Container is sealed")
	}

	reflectionValue := reflect.ValueOf(value)
	reflectionType := reflectionValue.Type()

	if !reflectionType.AssignableTo(interfaceType) {
		panic(fmt.Sprintf("Implementation type %s is not assignable to %s", reflectionType, interfaceType))
	}

	c.dependencies[interfaceType] = DependencyTarget{
		ReflectionValue: reflectionValue,
		Value:           value,
	}
}

func (c *ContainerImpl) SetValuesFromFunc(fns ...any) error {
	if c.isSealed {
		panic("Container is sealed")
	}

	for _, fn := range fns {
		fnValue := reflect.ValueOf(fn)

		fnType := fnValue.Type()
		if fnType.Kind() != reflect.Func {
			panic("Provider must be a function")
		}

		outNum := fnType.NumOut()
		lastOutTypeIsError := outNum > 0 && fnType.Out(outNum-1).AssignableTo(reflectTypeError)

		var interfaceCount int
		if lastOutTypeIsError {
			interfaceCount = outNum - 1
		} else {
			interfaceCount = outNum
		}

		inArgs := c.resolveInArgs(fnType)
		outArgs := fnValue.Call(inArgs)

		if len(outArgs) == 0 {
			continue
		}

		if lastOutTypeIsError && !outArgs[outNum-1].IsNil() {
			return outArgs[outNum-1].Interface().(error) //nolint:forcetypeassert
		}

		for i := range interfaceCount {
			interfaceType := fnType.Out(i)
			reflectionValue := outArgs[i]

			c.dependencies[interfaceType] = DependencyTarget{
				ReflectionValue: reflectionValue,
				Value:           reflectionValue.Interface(),
			}
		}
	}

	return nil
}

func (c *ContainerImpl) Seal() {
	c.isSealed = true
}

func (c *ContainerImpl) Resolve(t reflect.Type) (DependencyTarget, bool) {
	if t.Implements(reflectTypeContainer) {
		return DependencyTarget{ReflectionValue: reflect.ValueOf(c), Value: c}, true
	}

	target, ok := c.dependencies[t]

	return target, ok
}

func (c *ContainerImpl) MustResolve(t reflect.Type) DependencyTarget {
	value, ok := c.Resolve(t)

	if !ok {
		panic("No implementation registered for type " + t.String())
	}

	return value
}

func (c *ContainerImpl) CreateLister(t reflect.Type) func() []DependencyTarget { //nolint:varnamelen
	if c.isSealed {
		panic("Container is sealed")
	}

	var implementations []DependencyTarget

	for it, target := range c.dependencies {
		if (t.Kind() == reflect.Interface && it.Implements(t)) || it.AssignableTo(t) {
			implementations = append(implementations, target)
		}
	}

	return func() []DependencyTarget {
		return implementations
	}
}

func (c *ContainerImpl) DynamicList(t reflect.Type) []DependencyTarget {
	var implementations []DependencyTarget

	for it, target := range c.dependencies {
		if (t.Kind() == reflect.Interface && it.Implements(t)) || it.AssignableTo(t) {
			implementations = append(implementations, target)
		}
	}

	return implementations
}

func (c *ContainerImpl) CreateInvoker(fn any) func() error {
	if c.isSealed {
		panic("Container is sealed")
	}

	fnValue := reflect.ValueOf(fn)

	fnType := fnValue.Type()
	if fnType.Kind() != reflect.Func {
		panic("CreateInvoker parameter must be a function")
	}

	inArgs := c.resolveInArgs(fnType)

	return func() error {
		outArgs := fnValue.Call(inArgs)

		if len(outArgs) > 0 && !outArgs[0].IsNil() {
			return outArgs[0].Interface().(error) //nolint:forcetypeassert
		}

		return nil
	}
}

func (c *ContainerImpl) DynamicInvoke(fn any) error {
	fnValue := reflect.ValueOf(fn)

	fnType := fnValue.Type()
	if fnType.Kind() != reflect.Func {
		panic("DynamicInvoke parameter must be a function")
	}

	inArgs := c.resolveInArgs(fnType)
	outArgs := fnValue.Call(inArgs)

	if len(outArgs) > 0 && !outArgs[0].IsNil() {
		return outArgs[0].Interface().(error) //nolint:forcetypeassert
	}

	return nil
}

func (c *ContainerImpl) resolveInArgs(fnType reflect.Type) []reflect.Value {
	numIn := fnType.NumIn()
	args := make([]reflect.Value, numIn)

	for i := range args {
		paramType := fnType.In(i)
		args[i] = c.MustResolve(paramType).ReflectionValue
	}

	return args
}
