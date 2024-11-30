package datafx

type ConfigDataSource struct {
	DSN string `conf:"DSN"`
}

type Config struct {
	Sources map[string]ConfigDataSource `conf:"SOURCES"`
}
