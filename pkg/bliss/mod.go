package bliss

import (
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"go.uber.org/fx"
)

var FxModule = fx.Module( //nolint:gochecknoglobals
	"bliss",
	logfx.FxModule,
	configfx.FxModule,
	httpfx.FxModule,
	datafx.FxModule,
)
