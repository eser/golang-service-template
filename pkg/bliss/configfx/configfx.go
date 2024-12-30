package configfx

import (
	"github.com/eser/go-service/pkg/bliss/di"
)

func RegisterDependencies(container di.Container) {
	cl := NewConfigManager()

	di.RegisterFor[ConfigLoader](container, cl)
}
