package bliss

import (
	"github.com/eser/go-service/pkg/bliss/configfx"
	"github.com/eser/go-service/pkg/bliss/datafx"
	"github.com/eser/go-service/pkg/bliss/httpfx"
	"github.com/eser/go-service/pkg/bliss/logfx"
	"go.uber.org/fx"
)

var Module = fx.Module( //nolint:gochecknoglobals
	"bliss",
	logfx.Module,
	configfx.Module,
	httpfx.Module,
	datafx.Module,
)
