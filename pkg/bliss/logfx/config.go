package logfx

type Config struct {
	Level     string `conf:"LEVEL"      default:"INFO"`
	AddSource bool   `conf:"ADD_SOURCE" default:"false"`
}
