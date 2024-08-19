package logfx

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

type (
	StackTrace  = []uintptr // []runtime.Frame
	StackTracer interface {
		StackTrace() StackTrace
	}
)

func ReplacerGenerator(prettyMode bool) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		if prettyMode {
			if attr.Key == slog.TimeKey || attr.Key == slog.LevelKey || attr.Key == slog.MessageKey {
				return slog.Attr{
					Key:   "",
					Value: slog.Value{},
				}
			}
		}

		if attr.Value.Kind() == slog.KindAny {
			if v, ok := attr.Value.Any().(error); ok {
				attr.Value = fmtErr(v)
			}
		}

		return attr
	}
}

// fmtErr returns a slog.GroupValue with keys "msg" and "trace". If the error
// does not implement interface { StackTrace() StackTrace }, the "trace"
// key is omitted.
func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr

	groupValues = append(groupValues, slog.String("msg", err.Error()))

	// Find the trace to the location of the first errors.New,
	// errors.Wrap, or errors.WithStack call.
	var stackTraceable StackTracer

	for err := err; err != nil; err = errors.Unwrap(err) {
		if x, ok := err.(StackTracer); ok {
			stackTraceable = x
		}
	}

	if stackTraceable != nil {
		groupValues = append(groupValues,
			slog.Any("trace", TraceLines(stackTraceable.StackTrace())),
		)
	}

	return slog.GroupValue(groupValues...)
}

func TraceLines(frames StackTrace) []string {
	traceLines := make([]string, len(frames))

	// Iterate in reverse to skip uninteresting, consecutive runtime frames at
	// the bottom of the trace.
	var (
		skippedCounter int
		skipping       bool = true
	)

	for i := len(frames) - 1; i >= 0; i-- {
		// Adapted from errors.Frame.MarshalText(), but avoiding repeated
		// calls to FuncForPC and FileLine.
		programCounter := frames[i] - 1
		functionAddress := runtime.FuncForPC(programCounter)

		if functionAddress == nil {
			traceLines[i] = "unknown"
			skipping = false

			continue
		}

		name := functionAddress.Name()

		if skipping && strings.HasPrefix(name, "runtime.") {
			skippedCounter++

			continue
		}

		skipping = false
		filename, lineNr := functionAddress.FileLine(programCounter)

		traceLines[i] = fmt.Sprintf("%s %s:%d", name, filename, lineNr)
	}

	return traceLines[:len(traceLines)-skippedCounter]
}
