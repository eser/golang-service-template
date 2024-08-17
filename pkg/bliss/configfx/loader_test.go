package configfx_test

import (
	"reflect"
	"testing"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Host string `conf:"host" default:"localhost"`
}

type TestConfigNested struct {
	TestConfig
	Port int `conf:"port" default:"8080"`
}

func TestLoadMeta(t *testing.T) {
	t.Parallel()

	t.Run("should get config meta", func(t *testing.T) {
		t.Parallel()

		config := TestConfig{} //nolint:exhaustruct

		cl := configfx.NewConfigLoader()
		meta, err := cl.LoadMeta(&config)

		expected := []configfx.ConfigItemMeta{
			{
				Name:            "host",
				Field:           meta.Children[0].Field,
				Type:            reflect.TypeFor[string](),
				IsRequired:      false,
				HasDefaultValue: true,
				DefaultValue:    "localhost",

				Children: nil,
			},
		}

		if assert.NoError(t, err) {
			assert.Equal(t, "root", meta.Name)
			assert.Nil(t, meta.Type)

			assert.ElementsMatch(t, expected, meta.Children)
		}
	})

	t.Run("should get config meta from nested definition", func(t *testing.T) {
		t.Parallel()

		config := TestConfigNested{} //nolint:exhaustruct

		cl := configfx.NewConfigLoader()
		meta, err := cl.LoadMeta(&config)

		expected := []configfx.ConfigItemMeta{
			{
				Name:            "host",
				Field:           meta.Children[0].Field,
				Type:            reflect.TypeFor[string](),
				IsRequired:      false,
				HasDefaultValue: true,
				DefaultValue:    "localhost",

				Children: nil,
			},
			{
				Name:            "port",
				Field:           meta.Children[1].Field,
				Type:            reflect.TypeFor[int](),
				IsRequired:      false,
				HasDefaultValue: true,
				DefaultValue:    "8080",

				Children: nil,
			},
		}

		if assert.NoError(t, err) {
			assert.Equal(t, "root", meta.Name)
			assert.Nil(t, meta.Type)

			assert.ElementsMatch(t, expected, meta.Children)
		}
	})
}
