package configfx_test

import (
	"os"
	"testing"

	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentEnv(t *testing.T) {
	t.Run("should return current environment", func(t *testing.T) {
		oldEnv, oldEnvOk := os.LookupEnv("ENV")
		defer func() {
			if oldEnvOk {
				os.Setenv("ENV", oldEnv)
			} else {
				os.Unsetenv("ENV")
			}
		}()

		t.Setenv("ENV", "production")

		expected := "production"
		actual := configfx.GetCurrentEnv()

		assert.Equal(t, expected, actual)
	})
}

func TestGetFilenamesForEnv(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("should populate .env files for development environment", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			".env",
			".env.development",
			".env.local",
			".env.development.local",
		}

		actual := configfx.GetFilenamesForEnv("development", ".env")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate .env files for test environment", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			".env",
			".env.test",
			".env.test.local",
		}

		actual := configfx.GetFilenamesForEnv("test", ".env")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate .env files from parent directory", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"../.env",
			"../.env.development",
			"../.env.local",
			"../.env.development.local",
		}

		actual := configfx.GetFilenamesForEnv("development", "../.env")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate .env files from sub directory", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"testdata/.env",
			"testdata/.env.development",
			"testdata/.env.local",
			"testdata/.env.development.local",
		}

		actual := configfx.GetFilenamesForEnv("development", "testdata/.env")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate json config files for development environment", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"config.json",
			"config.development.json",
			"config.local.json",
			"config.development.local.json",
		}

		actual := configfx.GetFilenamesForEnv("development", "config.json")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate json config files for test environment", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"config.json",
			"config.test.json",
			"config.test.local.json",
		}

		actual := configfx.GetFilenamesForEnv("test", "config.json")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate json config files from parent directory", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"../config.json",
			"../config.development.json",
			"../config.local.json",
			"../config.development.local.json",
		}

		actual := configfx.GetFilenamesForEnv("development", "../config.json")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate json config files from sub directory", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"testdata/config.json",
			"testdata/config.development.json",
			"testdata/config.local.json",
			"testdata/config.development.local.json",
		}

		actual := configfx.GetFilenamesForEnv("development", "testdata/config.json")

		assert.ElementsMatch(t, expected, actual)
	})
}
