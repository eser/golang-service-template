package configfx

import (
	"fmt"

	"github.com/eser/go-service/pkg/bliss/configfx/envparser"
	"github.com/eser/go-service/pkg/bliss/configfx/jsonparser"
	"github.com/eser/go-service/pkg/bliss/lib"
)

func (dcl *ConfigLoaderImpl) FromEnvFileSingle(filename string) ConfigResource {
	return func(target *map[string]any) error {
		err := envparser.TryParseFiles(target, filename)
		if err != nil {
			return fmt.Errorf("failed to parse env file: %w", err)
		}

		return nil
	}
}

func (dcl *ConfigLoaderImpl) FromEnvFile(filename string) ConfigResource {
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

func (dcl *ConfigLoaderImpl) FromSystemEnv() ConfigResource {
	return func(target *map[string]any) error {
		lib.EnvOverrideVariables(target)

		return nil
	}
}

func (dcl *ConfigLoaderImpl) FromJsonFileSingle(filename string) ConfigResource {
	return func(target *map[string]any) error {
		err := jsonparser.TryParseFiles(target, filename)
		if err != nil {
			return fmt.Errorf("failed to parse json file: %w", err)
		}

		return nil
	}
}

func (dcl *ConfigLoaderImpl) FromJsonFile(filename string) ConfigResource {
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
