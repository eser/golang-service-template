package datafx

import "strings"

type Dialect string

const (
	// DialectPostgresPgx Dialect = "pgx".
	DialectPostgres Dialect = "postgres"
	DialectSQLite   Dialect = "sqlite"
	DialectMySQL    Dialect = "mysql"
)

func DetermineDialect(dsn string) Dialect {
	dsnLower := strings.ToLower(dsn)

	// if strings.HasPrefix(dsnLower, "pgx://") {
	// 	return DialectPostgresPgx
	// }

	if strings.HasPrefix(dsnLower, "postgres://") {
		return DialectPostgres
	}

	if strings.HasPrefix(dsnLower, "mysql://") {
		return DialectMySQL
	}

	if strings.HasPrefix(dsnLower, "sqlite://") {
		return DialectSQLite
	}

	// Default to postgres if cannot determine
	return DialectPostgres
}
