package jsonparser_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/configfx/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestTryParseFiles(t *testing.T) {
	t.Parallel()

	t.Run("should parse a json config file", func(t *testing.T) {
		t.Parallel()

		m := map[string]any{}
		err := jsonparser.TryParseFiles(&m, "./testdata/config.json")

		if assert.NoError(t, err) {
			assert.Equal(t, "env", m["TEST"])
			assert.Equal(t, "env!", m["TEST2__TEST3"])
		}
	})

	t.Run("should parse multiple json config files", func(t *testing.T) {
		t.Parallel()

		m := map[string]any{}
		err := jsonparser.TryParseFiles(&m, "./testdata/config.json", "./testdata/config.development.json")

		if assert.NoError(t, err) {
			assert.Equal(t, "env-development", m["TEST"])
			assert.Equal(t, "env!!", m["TEST2__TEST3"])
		}
	})
}
