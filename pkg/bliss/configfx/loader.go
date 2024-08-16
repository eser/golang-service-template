package configfx

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	tagConf    = "conf"
	tagDefault = "default"
)

var ErrNotStruct = errors.New("not a struct")

type ConfigItemMeta struct {
	Name            string
	Type            reflect.Type
	HasDefaultValue bool
	DefaultValue    string

	Children []ConfigItemMeta
}

type ConfigResource func(target *map[string]any) error

type ConfigLoader interface {
	LoadMeta(i any) (ConfigItemMeta, error)
	LoadMap(resources ...ConfigResource) (*map[string]any, error)
	Load(i any, resources ...ConfigResource) error

	FromEnvFileSingle(filename string) ConfigResource
	FromEnvFile(filename string) ConfigResource
	FromSystemEnv() ConfigResource

	FromJsonFileSingle(filename string) ConfigResource
	FromJsonFile(filename string) ConfigResource
}

type ConfigLoaderImpl struct{}

func NewConfigLoader() ConfigLoader { //nolint:ireturn
	return &ConfigLoaderImpl{}
}

func (dcl *ConfigLoaderImpl) LoadMeta(i any) (ConfigItemMeta, error) {
	r := reflect.ValueOf(i).Elem()

	children, err := reflectMeta(r)
	if err != nil {
		return ConfigItemMeta{}, err
	}

	return ConfigItemMeta{
		Name:            "root",
		Type:            nil,
		HasDefaultValue: false,
		DefaultValue:    "",

		Children: children,
	}, nil
}

func (dcl *ConfigLoaderImpl) LoadMap(resources ...ConfigResource) (*map[string]any, error) {
	target := map[string]any{}

	for _, resource := range resources {
		err := resource(&target)
		if err != nil {
			return nil, err
		}
	}

	return &target, nil
}

func (dcl *ConfigLoaderImpl) Load(i any, resources ...ConfigResource) error {
	meta, err := dcl.LoadMeta(i)
	if err != nil {
		return err
	}

	target, err := dcl.LoadMap(resources...)
	if err != nil {
		return err
	}

	for _, child := range meta.Children {
		fmt.Println(child)
	}
	fmt.Println(target)

	return nil
}

func reflectMeta(r reflect.Value) ([]ConfigItemMeta, error) {
	result := []ConfigItemMeta{}

	if r.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	for i := range r.NumField() {
		fieldType := r.Type().Field(i)

		if fieldType.Anonymous {
			children, err := reflectMeta(r.Field(i))
			if err != nil {
				return nil, err
			}

			if children != nil {
				result = append(result, children...)
			}

			continue
		}

		tag, hasTag := fieldType.Tag.Lookup(tagConf)
		if !hasTag {
			continue
		}

		defaultValue, hasDefaultValue := fieldType.Tag.Lookup(tagDefault)

		var children []ConfigItemMeta = nil

		if fieldType.Type.Kind() == reflect.Struct {
			var err error

			children, err = reflectMeta(r.Field(i))
			if err != nil {
				return nil, err
			}
		}

		result = append(result, ConfigItemMeta{
			Name:            tag,
			Type:            fieldType.Type,
			HasDefaultValue: hasDefaultValue,
			DefaultValue:    defaultValue,

			Children: children,
		})
	}

	return result, nil
}
