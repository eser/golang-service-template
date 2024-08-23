package app

import (
	"github.com/eser/go-service/pkg/bliss"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"go.uber.org/fx"
)

func Run() {
	app := fx.New(
		fx.WithLogger(logfx.GetFxLogger),
		bliss.FxModule,
		FxModule,
	)

	app.Run()
}
