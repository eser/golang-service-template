package results

import (
	"fmt"
	"log/slog"

	"github.com/eser/go-service/pkg/bliss/lib"
)

type ResultOccurrence struct {
	Definition *ResultDef

	InnerResult Result
	Error       error
	Attributes  []slog.Attr
}

var _ Result = (*ResultOccurrence)(nil)

func (r ResultOccurrence) IsOk() bool {
	return r.Error == nil
}

func (r ResultOccurrence) String() string {
	attrsStr := lib.SerializeSlogAttrs(r.AllAttributes()...)

	if r.InnerResult != nil {
		return fmt.Sprintf("[%s] %s (%s): %s", r.Definition.Code, r.Definition.Message, attrsStr, r.InnerResult.String())
	}

	return fmt.Sprintf("[%s] %s (%s)", r.Definition.Code, r.Definition.Message, attrsStr)
}

func (r ResultOccurrence) AllAttributes() []slog.Attr {
	return append(r.Definition.Attributes, r.Attributes...)
}

func (r ResultOccurrence) Unwrap() Result {
	return r.InnerResult
}

func (r ResultOccurrence) WithAttribute(logAttrs ...slog.Attr) ResultOccurrence {
	r.Attributes = append(r.Attributes, logAttrs...)

	return r
}

func (r ResultOccurrence) WrapResult(result Result) ResultOccurrence {
	r.InnerResult = result

	return r
}
