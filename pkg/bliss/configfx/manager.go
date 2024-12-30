package configfx

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/eser/go-service/pkg/bliss/results"
)

var ErrNotStruct = results.Define("ERRBC00001", "not a struct") //nolint:gochecknoglobals

type ConfigManager struct{}

var _ ConfigLoader = (*ConfigManager)(nil)

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

func (cl *ConfigManager) LoadMeta(i any) (ConfigItemMeta, error) {
	r := reflect.ValueOf(i).Elem() //nolint:varnamelen

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

// ------------------------
// Load Methods
// ------------------------

func (cl *ConfigManager) LoadMap(resources ...ConfigResource) (*map[string]any, error) {
	target := make(map[string]any)

	for _, resource := range resources {
		err := resource(&target)
		if err != nil {
			return nil, err
		}
	}

	return &target, nil
}

func (cl *ConfigManager) Load(i any, resources ...ConfigResource) error {
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

func (cl *ConfigManager) LoadDefaults(i any) error {
	return cl.Load(
		i,
		cl.FromJsonFile("config.json"),
		cl.FromEnvFile(".env"),
		cl.FromSystemEnv(),
	)
}

func reflectMeta(r reflect.Value) ([]ConfigItemMeta, error) { //nolint:varnamelen
	result := make([]ConfigItemMeta, 0)

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

		tag, hasTag := structFieldType.Tag.Lookup(TagConf)
		if !hasTag {
			continue
		}

		_, isRequired := structFieldType.Tag.Lookup(TagRequired)
		defaultValue, hasDefaultValue := structFieldType.Tag.Lookup(TagDefault)

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

func reflectSet(meta ConfigItemMeta, prefix string, target *map[string]any) { //nolint:funlen,cyclop
	for _, child := range meta.Children {
		key := prefix + child.Name

		if child.Type.Kind() == reflect.Map {
			// Create a new map
			newMap := reflect.MakeMap(child.Type)

			// Find all keys that start with our prefix
			prefix := key + Separator
			for key := range *target {
				if !strings.HasPrefix(key, prefix) {
					continue
				}

				// Extract the map key from the flattened key
				mapKey := strings.TrimPrefix(key, prefix)
				if idx := strings.Index(mapKey, Separator); idx != -1 {
					mapKey = mapKey[:idx]
				}

				// Create and set the map value
				valueType := child.Type.Elem()
				mapValue := reflect.New(valueType).Elem()

				// Recursively set the fields of the map value
				subMeta := ConfigItemMeta{
					Name:            mapKey,
					Field:           mapValue,
					Type:            valueType,
					IsRequired:      child.IsRequired,
					HasDefaultValue: child.HasDefaultValue,
					DefaultValue:    child.DefaultValue,

					Children: nil,
				}

				if valueType.Kind() == reflect.Struct {
					children, _ := reflectMeta(mapValue)
					subMeta.Children = children
				}

				reflectSet(subMeta, prefix+mapKey+Separator, target)

				// Set the value in the map
				newMap.SetMapIndex(reflect.ValueOf(mapKey), mapValue)
			}

			child.Field.Set(newMap)

			continue
		}

		if child.Type.Kind() == reflect.Struct {
			reflectSet(child, key+Separator, target)

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

func reflectSetField(field reflect.Value, fieldType reflect.Type, value string) { //nolint:funlen,cyclop
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
