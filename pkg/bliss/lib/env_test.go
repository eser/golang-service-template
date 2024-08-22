package lib_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/stretchr/testify/assert"
)

func TestEnvAwareFilenames(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("should populate .env files for development environment", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			".env",
			".env.development",
			".env.local",
			".env.development.local",
		}

		actual := lib.EnvAwareFilenames("development", ".env")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate .env files for test environment", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			".env",
			".env.test",
			".env.test.local",
		}

		actual := lib.EnvAwareFilenames("test", ".env")

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

		actual := lib.EnvAwareFilenames("development", "../.env")

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

		actual := lib.EnvAwareFilenames("development", "testdata/.env")

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

		actual := lib.EnvAwareFilenames("development", "config.json")

		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("should populate json config files for test environment", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"config.json",
			"config.test.json",
			"config.test.local.json",
		}

		actual := lib.EnvAwareFilenames("test", "config.json")

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

		actual := lib.EnvAwareFilenames("development", "../config.json")

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

		actual := lib.EnvAwareFilenames("development", "testdata/config.json")

		assert.ElementsMatch(t, expected, actual)
	})
}

func TestEnvGetCurrent(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "should return development when ENV is empty",
			envValue: "",
			expected: "development",
		},
		{
			name:     "should return lowercase value of ENV",
			envValue: "PRODUCTION",
			expected: "production",
		},
		{
			name:     "should return lowercase value of ENV with leading/trailing spaces",
			envValue: "STAGING",
			expected: "staging",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ENV", tt.envValue)

			actual := lib.EnvGetCurrent()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestEnvOverrideVariables(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		expectedArgs map[string]any
	}{
		{
			name: "should override variables",
			env:  map[string]string{"ENV": "development"},
			expectedArgs: map[string]any{
				"ENV": "development",
			},
		},
		{
			name: "should override multiple variables",
			env:  map[string]string{"ENV": "development", "DEBUG": "true", "PORT": "8080"},
			expectedArgs: map[string]any{
				"ENV":   "development",
				"DEBUG": "true",
				"PORT":  "8080",
			},
		},
		{
			name:         "should handle empty environment",
			env:          map[string]string{},
			expectedArgs: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := make(map[string]any)

			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			lib.EnvOverrideVariables(&m)

			for k, v := range tt.expectedArgs {
				assert.Equal(t, v, m[k])
			}
		})
	}
}
