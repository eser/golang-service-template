package datafx

import "strings"

type Dialect string

const (
	DialectPostgres Dialect = "postgres"
	DialectMySQL    Dialect = "mysql"
)

func DetermineDialect(dsn string) Dialect {
	dsnLower := strings.ToLower(dsn)

	if strings.HasPrefix(dsnLower, "postgres://") {
		return DialectPostgres
	}

	if strings.HasPrefix(dsnLower, "mysql://") {
		return DialectMySQL
	}

	// Default to postgres if cannot determine
	return DialectPostgres
}
