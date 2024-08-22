package results

import (
	"log/slog"
)

type Result interface {
	IsOk() bool

	String() string
	AllAttributes() []slog.Attr

	Unwrap() Result
}
