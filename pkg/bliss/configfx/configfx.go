package configfx

import (
	"github.com/eser/go-service/pkg/bliss/di"
)

func RegisterDependencies(container di.Container) {
	cl := NewConfigLoader()

	di.RegisterFor[ConfigLoader](container, cl)
}
