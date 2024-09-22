package envparser_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/configfx/envparser"
	"github.com/stretchr/testify/assert"
)

func TestTryParseFiles(t *testing.T) {
	t.Parallel()

	t.Run("should parse a .env file", func(t *testing.T) {
		t.Parallel()

		m := map[string]any{}
		err := envparser.TryParseFiles(&m, "./testdata/.env")

		if assert.NoError(t, err) {
			assert.Equal(t, "env", m["TEST"])
		}
	})

	t.Run("should parse multiple .env files", func(t *testing.T) {
		t.Parallel()

		m := map[string]any{}
		err := envparser.TryParseFiles(&m, "./testdata/.env", "./testdata/.env.development")

		if assert.NoError(t, err) {
			assert.Equal(t, "env-development", m["TEST"])
		}
	})
}
