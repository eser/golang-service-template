package configfx

import "github.com/eser/go-service/pkg/bliss/configfx/envparser"

func (dcl *ConfigLoaderImpl) FromEnvFileSingle(filenames ...string) ConfigResource {
	return func(meta ConfigItemMeta, target *map[string]any) error {
		err := envparser.TryParseFiles(target, filenames...)

		return err
	}
}

func (dcl *ConfigLoaderImpl) FromEnvFile(filenames ...string) ConfigResource {
	return func(meta ConfigItemMeta, target *map[string]any) error {
		// TODO(@eser): Implement this function

		return nil
	}
}

func (dcl *ConfigLoaderImpl) FromSystemEnv() ConfigResource {
	return func(meta ConfigItemMeta, target *map[string]any) error {
		// TODO(@eser): Implement this function

		return nil
	}
}

func (dcl *ConfigLoaderImpl) FromJsonFileSingle(filenames ...string) ConfigResource {
	return func(meta ConfigItemMeta, target *map[string]any) error {
		// TODO(@eser): Implement this function

		return nil
	}
}

func (dcl *ConfigLoaderImpl) FromJsonFile(filenames ...string) ConfigResource {
	return func(meta ConfigItemMeta, target *map[string]any) error {
		// TODO(@eser): Implement this function

		return nil
	}
}
