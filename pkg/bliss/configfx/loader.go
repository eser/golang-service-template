package configfx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/eser/go-service/pkg/bliss/configfx/envparser"
	"github.com/eser/go-service/pkg/bliss/configfx/jsonparser"
	"github.com/eser/go-service/pkg/bliss/lib"
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

type ConfigResource func(meta ConfigItemMeta, target *map[string]any) error

type ConfigLoader interface {
	LoadMeta(i any) (ConfigItemMeta, error)
	Load(i any, resources ...ConfigResource) error

	FromEnvFileSingle(filenames ...string) ConfigResource
	FromEnvFile(filenames ...string) ConfigResource
	FromSystemEnv() ConfigResource

	FromJsonFileSingle(filenames ...string) ConfigResource
	FromJsonFile(filenames ...string) ConfigResource
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
		Name:     "root",
		Children: children,
	}, nil
}

func (dcl *ConfigLoaderImpl) Load(i any, resources ...ConfigResource) error {
	meta, err := dcl.LoadMeta(i)
	if err != nil {
		return err
	}

	for _, child := range meta.Children {
		fmt.Println(child)
	}

	return nil
}

func tryLoadEnv(m *map[string]any) error {
	env := lib.EnvGetCurrent()
	filenames := lib.EnvAwareFilenames(env, ".env")

	err := envparser.TryParseFiles(m, filenames...)
	if err != nil {
		return err //nolint:wrapcheck
	}

	lib.EnvOverrideVariables(m)

	return nil
}

func tryLoadJson(m *map[string]any) error {
	env := lib.EnvGetCurrent()
	filenames := lib.EnvAwareFilenames(env, "config.json")

	err := jsonparser.TryParseFiles(m, filenames...)
	if err != nil {
		return err //nolint:wrapcheck
	}

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
