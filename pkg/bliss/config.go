package bliss

type BaseConfig struct {
	Env string `conf:"env" default:"development"`

	// AppName           string `conf:"APP_NAME" default:"go-service"`
	// LogTarget         string `conf:"LOG_TARGET" default:"stdout"`
	// Port              string `conf:"PORT" default:"8080"`
	// JwtSignature      string `conf:"JWT_SIGNATURE"`
	// CorsOrigin        string `conf:"CORS_ORIGIN"`
	// CorsStrictHeaders bool   `conf:"CORS_STRICT_HEADERS"`
	// DataConnstr       string `conf:"DATA_CONNSTR"`
}
