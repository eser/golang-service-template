package bliss

import (
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"go.uber.org/fx"
)

//nolint:gochecknoglobals
var Module = fx.Module(
	"bliss",
	configfx.Module,
	httpfx.Module,
)
