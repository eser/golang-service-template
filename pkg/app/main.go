package app

import (
	"github.com/eser/go-service/pkg/bliss"
	"go.uber.org/fx"
)

func Run() {
	app := fx.New(
		// fx.WithLogger(bliss.GetFxLogger),
		bliss.Module,
		Module,
	)

	app.Run()
}
