package logfx

type Config struct {
	Level      string `conf:"LEVEL"      default:"INFO"`
	PrettyMode bool   `conf:"PRETTY"     default:"true"`
	AddSource  bool   `conf:"ADD_SOURCE" default:"false"`
}
