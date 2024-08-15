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
	tagConf = "conf"
)

var ErrNotStruct = errors.New("not a struct")

type ConfigLoader struct{}

type ConfigMeta struct {
	Name     string
	Type     reflect.Type
	Children *[]ConfigMeta
}

func (dcl *ConfigLoader) TryLoadEnv(m *map[string]string) error {
	env := lib.EnvGetCurrent()
	filenames := lib.EnvAwareFilenames(env, ".env")

	err := envparser.TryParseFiles(m, filenames...)
	if err != nil {
		return err //nolint:wrapcheck
	}

	lib.EnvOverrideVariables(m)

	return nil
}

func (dcl *ConfigLoader) TryLoadJson(m *map[string]any) error {
	env := lib.EnvGetCurrent()
	filenames := lib.EnvAwareFilenames(env, "config.json")

	err := jsonparser.TryParseFiles(m, filenames...)
	if err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

func reflectMeta(r reflect.Value) (*[]ConfigMeta, error) {
	result := []ConfigMeta{}

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

			result = append(result, *children...)

			continue
		}

		tag, ok := fieldType.Tag.Lookup(tagConf)
		if !ok {
			continue
		}

		result = append(result, ConfigMeta{
			Name: tag,
			Type: fieldType.Type,
		})
	}

	return &result, nil
}

func (dcl *ConfigLoader) LoadMeta(i any) (*ConfigMeta, error) {
	r := reflect.ValueOf(i).Elem()
	children, err := reflectMeta(r)
	if err != nil {
		return nil, err
	}

	return &ConfigMeta{
		Name:     "root",
		Children: children,
	}, nil
}

func (dcl *ConfigLoader) Load(i any) error {
	meta, err := dcl.LoadMeta(i)
	if err != nil {
		return err
	}

	for _, child := range *meta.Children {
		fmt.Printf("fieldName: %s fieldType: %s\n", child.Name, child.Type.Name())
	}

	return nil
}
