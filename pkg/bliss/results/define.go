package results

import (
	"log/slog"
	"strings"
)

type Definition struct {
	Code    string
	Message string

	Attributes []slog.Attr
}

func Define(code string, message string, attributes ...slog.Attr) *Definition {
	if attributes == nil {
		attributes = []slog.Attr{}
	}

	return &Definition{
		Code:    code,
		Message: message,

		Attributes: attributes,
	}
}

func (r *Definition) New(messages ...string) ResultImpl {
	if messages == nil {
		messages = []string{}
	}

	return ResultImpl{
		Definition: r,

		InnerError:      nil,
		InnerMessage:    strings.Join(messages, ", "),
		InnerAttributes: []slog.Attr{},
	}
}

func (r *Definition) Wrap(err error, messages ...string) ResultImpl {
	if messages == nil {
		messages = []string{}
	}

	return ResultImpl{
		Definition: r,

		InnerError:      err,
		InnerMessage:    strings.Join(messages, ", "),
		InnerAttributes: []slog.Attr{},
	}
}
