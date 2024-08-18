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

func replaceAttr(groups []string, attr slog.Attr) slog.Attr {
	switch attr.Value.Kind() { //nolint:gocritic,wsl,exhaustive
	// other cases

	case slog.KindAny:
		switch v := attr.Value.Any().(type) { //nolint:gocritic
		case error:
			attr.Value = fmtErr(v)
		}
	}

	return attr
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
			slog.Any("trace", traceLines(stackTraceable.StackTrace())),
		)
	}

	return slog.GroupValue(groupValues...)
}

func traceLines(frames StackTrace) []string {
	traceLines := make([]string, len(frames))

	// Iterate in reverse to skip uninteresting, consecutive runtime frames at
	// the bottom of the trace.
	var (
		skippedCounter int
		skipping       bool = true
	)

	for i := len(frames) - 1; i >= 0; i-- { //nolint:varnamelen
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
