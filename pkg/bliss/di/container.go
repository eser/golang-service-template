package di

import (
	"fmt"
	"reflect"
)

type Provider func(args []any) any

type DependencyTarget struct {
	Value           any
	ReflectionValue reflect.Value
}

// Container defines the methods for dependency injection container.
type Container struct {
	dependencies map[reflect.Type]DependencyTarget

	isSealed bool
}

var reflectTypeError = reflect.TypeFor[error]() //nolint:gochecknoglobals

// NewContainer creates a new dependency injection container.
func NewContainer() *Container {
	return &Container{
		isSealed: false,

		dependencies: make(map[reflect.Type]DependencyTarget),
	}
}

func (c *Container) SetValue(value any) {
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

func (c *Container) SetValueFor(interfaceType reflect.Type, value any) {
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

func (c *Container) SetValuesFromFunc(fns ...any) error {
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

func (c *Container) Seal() {
	c.isSealed = true
}

func (c *Container) Resolve(t reflect.Type) (DependencyTarget, bool) {
	target, isOk := c.dependencies[t]

	if !isOk {
		for _, target := range c.dependencies {
			if target.ReflectionValue.Type().Implements(t) {
				return target, true
			}
		}
	}

	return target, isOk
}

func (c *Container) MustResolve(t reflect.Type) DependencyTarget {
	value, ok := c.Resolve(t)

	if !ok {
		panic("No implementation registered for type " + t.String())
	}

	return value
}

func (c *Container) CreateLister(t reflect.Type) func() []DependencyTarget { //nolint:varnamelen
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

func (c *Container) DynamicList(t reflect.Type) []DependencyTarget {
	var implementations []DependencyTarget

	for it, target := range c.dependencies {
		if (t.Kind() == reflect.Interface && it.Implements(t)) || it.AssignableTo(t) {
			implementations = append(implementations, target)
		}
	}

	return implementations
}

func (c *Container) CreateInvoker(fn any) func() error {
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

func (c *Container) DynamicInvoke(fn any) error {
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

func (c *Container) resolveInArgs(fnType reflect.Type) []reflect.Value {
	numIn := fnType.NumIn()
	args := make([]reflect.Value, numIn)

	for i := range args {
		paramType := fnType.In(i)
		args[i] = c.MustResolve(paramType).ReflectionValue
	}

	return args
}
