package logfx

import (
	"log/slog"

	"go.uber.org/fx/fxevent"
)

type (
	FxLogger struct {
		*slog.Logger
	}
)

func GetFxLogger(logger *slog.Logger) fxevent.Logger { //nolint:ireturn
	return &FxLogger{logger}
}

func (l FxLogger) LogEvent(event fxevent.Event) { //nolint:cyclop
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logOnStartExecuting(e)
	case *fxevent.OnStartExecuted:
		l.logOnStartExecuted(e)
	case *fxevent.OnStopExecuting:
		l.logOnStopExecuting(e)
	case *fxevent.OnStopExecuted:
		l.logOnStopExecuted(e)
	case *fxevent.Supplied:
		l.logSupplied(e)
	case *fxevent.Provided:
		l.logProvided(e)
	case *fxevent.Decorated:
		l.logDecorated(e)
	case *fxevent.Invoking:
		l.logInvoking(e)
	case *fxevent.Started:
		l.logStarted(e)
	case *fxevent.LoggerInitialized:
		l.logLoggerInitialized(e)
	}
}

func (l *FxLogger) logOnStartExecuting(e *fxevent.OnStartExecuting) {
	l.Logger.Debug(
		"OnStart hook executing: ",
		slog.String("callee", e.FunctionName),
		slog.String("caller", e.CallerName),
	)
}

func (l *FxLogger) logOnStartExecuted(e *fxevent.OnStartExecuted) {
	if e.Err != nil {
		l.Logger.Debug(
			"OnStart hook failed: ",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
			slog.Any("error", e.Err),
		)

		return
	}

	l.Logger.Debug(
		"OnStart hook executed: ",
		slog.String("callee", e.FunctionName),
		slog.String("caller", e.CallerName),
		slog.String("runtime", e.Runtime.String()),
	)
}

func (l *FxLogger) logOnStopExecuting(e *fxevent.OnStopExecuting) {
	l.Logger.Debug(
		"OnStop hook executing: ",
		slog.String("callee", e.FunctionName),
		slog.String("caller", e.CallerName),
	)
}

func (l *FxLogger) logOnStopExecuted(e *fxevent.OnStopExecuted) {
	if e.Err != nil {
		l.Logger.Debug(
			"OnStop hook failed: ",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
			slog.Any("error", e.Err),
		)

		return
	}

	l.Logger.Debug(
		"OnStop hook executed: ",
		slog.String("callee", e.FunctionName),
		slog.String("caller", e.CallerName),
		slog.String("runtime", e.Runtime.String()),
	)
}

func (l *FxLogger) logSupplied(e *fxevent.Supplied) {
	l.Logger.Debug(
		"supplied: ",
		slog.String("type", e.TypeName),
		slog.Any("error", e.Err),
	)
}

func (l *FxLogger) logProvided(e *fxevent.Provided) {
	for _, rtype := range e.OutputTypeNames {
		l.Logger.Debug(
			"provided: ",
			slog.String("constructor", e.ConstructorName),
			slog.String("type", rtype),
		)
	}
}

func (l *FxLogger) logDecorated(e *fxevent.Decorated) {
	for _, rtype := range e.OutputTypeNames {
		l.Logger.Debug("decorated: ",
			slog.String("decorator", e.DecoratorName),
			slog.String("type", rtype),
		)
	}
}

func (l *FxLogger) logInvoking(e *fxevent.Invoking) {
	l.Logger.Debug(
		"invoking: ",
		slog.String("function", e.FunctionName),
	)
}

func (l *FxLogger) logStarted(e *fxevent.Started) {
	if e.Err == nil {
		l.Logger.Debug("started")
	}
}

func (l *FxLogger) logLoggerInitialized(e *fxevent.LoggerInitialized) {
	if e.Err == nil {
		l.Logger.Debug(
			"initialized: custom fxevent.Logger",
			slog.String("constructor", e.ConstructorName),
		)
	}
}
