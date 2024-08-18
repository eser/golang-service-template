package bliss

import (
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"go.uber.org/fx"
)

//nolint:gochecknoglobals
var Module = fx.Module(
	"bliss",
	// fx.WithLogger(logfx.GetFxLogger),
	logfx.Module,
	configfx.Module,
	httpfx.Module,
)
