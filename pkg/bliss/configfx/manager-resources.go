package configfx

import (
	"fmt"

	"github.com/eser/go-service/pkg/bliss/configfx/envparser"
	"github.com/eser/go-service/pkg/bliss/configfx/jsonparser"
	"github.com/eser/go-service/pkg/bliss/lib"
)

func (cl *ConfigManager) FromEnvFileDirect(filename string) ConfigResource {
	return func(target *map[string]any) error {
		err := envparser.TryParseFiles(target, filename)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		return nil
	}
}

func (cl *ConfigManager) FromEnvFile(filename string) ConfigResource {
	return func(target *map[string]any) error {
		env := lib.EnvGetCurrent()
		filenames := lib.EnvAwareFilenames(env, filename)

		err := envparser.TryParseFiles(target, filenames...)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		return nil
	}
}

func (cl *ConfigManager) FromSystemEnv() ConfigResource {
	return func(target *map[string]any) error {
		lib.EnvOverrideVariables(target)

		return nil
	}
}

func (cl *ConfigManager) FromJsonFileDirect(filename string) ConfigResource {
	return func(target *map[string]any) error {
		err := jsonparser.TryParseFiles(target, filename)
		if err != nil {
			return fmt.Errorf("failed to parse json file: %w", err)
		}

		return nil
	}
}

func (cl *ConfigManager) FromJsonFile(filename string) ConfigResource {
	return func(target *map[string]any) error {
		env := lib.EnvGetCurrent()
		filenames := lib.EnvAwareFilenames(env, filename)

		err := jsonparser.TryParseFiles(target, filenames...)
		if err != nil {
			return fmt.Errorf("failed to parse json file: %w", err)
		}

		return nil
	}
}
