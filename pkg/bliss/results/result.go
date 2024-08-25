package results

import (
	"log/slog"
	"strings"

	"github.com/eser/go-service/pkg/bliss/lib"
)

type Result interface {
	error
	Unwrap() error

	IsError() bool
	String() string
	Attributes() []slog.Attr
}

type ResultImpl struct { //nolint:errname
	Definition *Definition

	InnerError      error
	InnerPayload    any
	InnerAttributes []slog.Attr
}

var _ Result = (*ResultImpl)(nil)

func (r ResultImpl) Error() string {
	return r.String()
}

func (r ResultImpl) Unwrap() error {
	return r.InnerError
}

func (r ResultImpl) IsError() bool {
	return r.InnerError != nil
}

func (r ResultImpl) String() string {
	attrs := r.Attributes()

	builder := strings.Builder{}
	builder.WriteRune('[')
	builder.WriteString(r.Definition.Code)
	builder.WriteString("] ")
	builder.WriteString(r.Definition.Message)

	if len(attrs) > 0 {
		builder.WriteString(" (")
		builder.WriteString(lib.SerializeSlogAttrs(attrs))
		builder.WriteRune(')')
	}

	if r.InnerError != nil {
		builder.WriteString(": ")
		builder.WriteString(r.InnerError.Error())
	}

	return builder.String()
}

func (r ResultImpl) Attributes() []slog.Attr {
	attrs := r.Definition.Attributes

	// if r.InnerPayload != nil {
	// 	attrs = append(attrs, slog.Any("payload", r.InnerPayload))
	// }

	return append(attrs, r.InnerAttributes...)
}

func (r ResultImpl) WithError(err error) ResultImpl {
	r.InnerError = err

	return r
}

func (r ResultImpl) WithAttribute(attributes ...slog.Attr) ResultImpl {
	r.InnerAttributes = append(r.InnerAttributes, attributes...)

	return r
}
