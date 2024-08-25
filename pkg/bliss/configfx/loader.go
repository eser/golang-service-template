package configfx

import (
	"reflect"
	"strconv"
	"time"

	"github.com/eser/go-service/pkg/bliss/results"
)

const (
	tagConf     = "conf"
	tagDefault  = "default"
	tagRequired = "required"
)

var ErrNotStruct = results.Define("ERRBC00001", "not a struct") //nolint:gochecknoglobals

type ConfigItemMeta struct {
	Name            string
	Field           reflect.Value
	Type            reflect.Type
	IsRequired      bool
	HasDefaultValue bool
	DefaultValue    string

	Children []ConfigItemMeta
}

type ConfigResource func(target *map[string]any) error

type ConfigLoader interface {
	LoadMeta(i any) (ConfigItemMeta, error)
	LoadMap(resources ...ConfigResource) (*map[string]any, error)
	Load(i any, resources ...ConfigResource) error

	FromEnvFileDirect(filename string) ConfigResource
	FromEnvFile(filename string) ConfigResource
	FromSystemEnv() ConfigResource

	FromJsonFileDirect(filename string) ConfigResource
	FromJsonFile(filename string) ConfigResource
}

type ConfigLoaderImpl struct{}

var _ ConfigLoader = (*ConfigLoaderImpl)(nil)

func NewConfigLoader() *ConfigLoaderImpl {
	return &ConfigLoaderImpl{}
}

func (cl *ConfigLoaderImpl) LoadMeta(i any) (ConfigItemMeta, error) {
	r := reflect.ValueOf(i).Elem()

	children, err := reflectMeta(r)
	if err != nil {
		return ConfigItemMeta{}, err
	}

	return ConfigItemMeta{
		Name:            "root",
		Field:           r,
		Type:            nil,
		IsRequired:      false,
		HasDefaultValue: false,
		DefaultValue:    "",

		Children: children,
	}, nil
}

func (cl *ConfigLoaderImpl) LoadMap(resources ...ConfigResource) (*map[string]any, error) {
	target := map[string]any{}

	for _, resource := range resources {
		err := resource(&target)
		if err != nil {
			return nil, err
		}
	}

	return &target, nil
}

func (cl *ConfigLoaderImpl) Load(i any, resources ...ConfigResource) error {
	meta, err := cl.LoadMeta(i)
	if err != nil {
		return err
	}

	target, err := cl.LoadMap(resources...)
	if err != nil {
		return err
	}

	reflectSet(meta, "", target)

	return nil
}

func reflectMeta(r reflect.Value) ([]ConfigItemMeta, error) {
	result := []ConfigItemMeta{}

	if r.Kind() != reflect.Struct {
		return nil, ErrNotStruct.New()
	}

	for i := range r.NumField() {
		structField := r.Field(i)
		structFieldType := r.Type().Field(i)

		if structFieldType.Anonymous {
			children, err := reflectMeta(structField)
			if err != nil {
				return nil, err
			}

			if children != nil {
				result = append(result, children...)
			}

			continue
		}

		tag, hasTag := structFieldType.Tag.Lookup(tagConf)
		if !hasTag {
			continue
		}

		_, isRequired := structFieldType.Tag.Lookup(tagRequired)
		defaultValue, hasDefaultValue := structFieldType.Tag.Lookup(tagDefault)

		var children []ConfigItemMeta = nil

		if structFieldType.Type.Kind() == reflect.Struct {
			var err error

			children, err = reflectMeta(structField)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, ConfigItemMeta{
			Name:            tag,
			Field:           structField,
			Type:            structFieldType.Type,
			IsRequired:      isRequired,
			HasDefaultValue: hasDefaultValue,
			DefaultValue:    defaultValue,

			Children: children,
		})
	}

	return result, nil
}

func reflectSet(meta ConfigItemMeta, prefix string, target *map[string]any) {
	for _, child := range meta.Children {
		key := prefix + child.Name

		if child.Type.Kind() == reflect.Struct {
			reflectSet(child, key+"__", target)

			continue
		}

		// Check if the target map has the key with the child name
		value, valueOk := (*target)[key].(string)
		if !valueOk {
			if child.HasDefaultValue {
				reflectSetField(child.Field, child.Type, child.DefaultValue)

				continue
			}

			if child.IsRequired {
				panic("missing required config value: " + child.Name)
			}

			continue
		}

		reflectSetField(child.Field, child.Type, value)
	}
}

func reflectSetField(field reflect.Value, fieldType reflect.Type, value string) {
	var finalValue reflect.Value

	switch fieldType {
	case reflect.TypeFor[string]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[int]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[int8]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[int16]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[int32]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[int64]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[uint]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[uint8]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[uint16]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[uint32]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[uint64]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[float32]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[float64]():
		finalValue = reflect.ValueOf(value)
	case reflect.TypeFor[bool]():
		boolValue, _ := strconv.ParseBool(value)
		finalValue = reflect.ValueOf(boolValue)
	case reflect.TypeFor[time.Duration]():
		durationValue, _ := time.ParseDuration(value)
		finalValue = reflect.ValueOf(durationValue)
	default:
		return
	}

	if field.Kind() == reflect.Ptr {
		// Handle pointer types by allocating a new instance
		ptr := reflect.New(fieldType.Elem())
		ptr.Elem().Set(finalValue)
		field.Set(ptr)

		return
	}

	// Set the field directly
	field.Set(finalValue)
}
