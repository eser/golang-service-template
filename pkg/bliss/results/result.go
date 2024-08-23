package results

import (
	"fmt"
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/lib"
)

type Result interface {
	error
	Unwrap() error

	IsError() bool
	String() string
	Attributes() []slog.Attr
}

type ResultImpl struct {
	Definition *Definition

	InnerError      error
	InnerMessage    string
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

	var attrsStr string
	if len(attrs) > 0 {
		attrsStr = fmt.Sprintf(" (%s)", lib.SerializeSlogAttrs(attrs))
	}

	if r.InnerError != nil {
		return fmt.Sprintf("[%s] %s%s: %s", r.Definition.Code, r.Definition.Message, attrsStr, r.InnerError.Error())
	}

	return fmt.Sprintf("[%s] %s%s", r.Definition.Code, r.Definition.Message, attrsStr)
}

func (r ResultImpl) Attributes() []slog.Attr {
	return append(r.Definition.Attributes, r.InnerAttributes...)
}

func (r ResultImpl) WithError(err error) ResultImpl {
	r.InnerError = err

	return r
}

func (r ResultImpl) WithAttribute(attributes ...slog.Attr) ResultImpl {
	r.InnerAttributes = append(r.InnerAttributes, attributes...)

	return r
}
