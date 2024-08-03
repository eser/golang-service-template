package bliss

import (
	"github.com/eser/go-service/pkg/bliss/config"
	"github.com/eser/go-service/pkg/bliss/httpservice"
	"go.uber.org/fx"
)

//nolint:gochecknoglobals
var Module = fx.Module(
	"bliss",
	config.Module,
	httpservice.Module,
)
