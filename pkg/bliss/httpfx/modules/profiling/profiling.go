package profiling

import (
	"net/http/pprof"

	"github.com/eser/go-service/pkg/bliss/httpfx"
)

func RegisterHttpRoutes(routes httpfx.Router, config *httpfx.Config) error {
	if !config.ProfilingEnabled {
		return nil
	}

	routes.GetMux().HandleFunc("/debug/pprof/", pprof.Index)
	routes.GetMux().HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	routes.GetMux().HandleFunc("/debug/pprof/profile", pprof.Profile)
	routes.GetMux().HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	routes.GetMux().HandleFunc("/debug/pprof/trace", pprof.Trace)

	return nil
}
