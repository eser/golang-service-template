package results

import (
	"log/slog"
)

type ResultDef struct {
	Code    string
	Message string

	Attributes []slog.Attr
}

func NewResultDef(code string, message string, logAttrs ...slog.Attr) *ResultDef {
	return &ResultDef{
		Code:    code,
		Message: message,

		Attributes: logAttrs,
	}
}

func (r *ResultDef) New() ResultOccurrence {
	return ResultOccurrence{
		Definition: r,

		InnerResult: nil,
		Error:       nil,
		Attributes:  []slog.Attr{},
	}
}

func (r *ResultDef) NewWithError(err error) ResultOccurrence {
	return ResultOccurrence{
		Definition: r,

		InnerResult: nil,
		Error:       err,
		Attributes:  []slog.Attr{},
	}
}

func (r *ResultDef) Wrap(result Result) ResultOccurrence {
	return ResultOccurrence{
		Definition: r,

		InnerResult: result,
		Error:       nil,
		Attributes:  []slog.Attr{},
	}
}

func (r *ResultDef) WrapWithError(result Result, err error) ResultOccurrence {
	return ResultOccurrence{
		Definition: r,

		InnerResult: result,
		Error:       err,
		Attributes:  []slog.Attr{},
	}
}
