package logfx_test

import (
	"bufio"
	"errors"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/eser/go-service/pkg/bliss/logfx"
	"go.uber.org/fx/fxevent"
)

type MockWriter struct {
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func generateFxLogger() (*logfx.FxLogger, *slog.Logger) {
	logger := slog.New(slog.NewJSONHandler(&MockWriter{}, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	return &logfx.FxLogger{Logger: logger}, logger
}

func TestGetFxLogger(t *testing.T) {
	t.Parallel()
	fxLogger := logfx.GetFxLogger(nil)
	if fxLogger == nil {
		t.Error("GetFxLogger() = nil, want not nil")
	}
}

func TestFxLogger_LogEvent(t *testing.T) {
	t.Parallel()
	fxLogger, _ := generateFxLogger()

	tests := []struct {
		name  string
		event fxevent.Event
		want  string
	}{
		{
			name: "OnStartExecuting",
			event: &fxevent.OnStartExecuting{
				FunctionName: "startFunc",
				CallerName:   "callerFunc",
			},
			want: `{"level":"debug","message":"OnStart hook executing: ","callee":"startFunc","caller":"callerFunc"}`,
		},
		{
			name: "OnStartExecuted",
			event: &fxevent.OnStartExecuted{
				FunctionName: "startFunc",
				CallerName:   "callerFunc",
				Runtime:      time.Second,
			},
			want: `{"level":"debug","message":"OnStart hook executing: ","callee":"startFunc","caller":"callerFunc"}`,
		},
		{
			name: "OnStartExecuted with err",
			event: &fxevent.OnStartExecuted{
				FunctionName: "startFunc",
				CallerName:   "callerFunc",
				Runtime:      time.Second,
				Err:          errors.New("error"),
			},
			want: `{"level":"debug","message":"OnStart hook failed: ","callee":"startFunc","caller":"callerFunc","error":"error"}`,
		},
		{
			name: "OnStopExecuting",
			event: &fxevent.OnStopExecuting{
				FunctionName: "stopFunc",
				CallerName:   "callerFunc",
			},
			want: `{"level":"debug","message":"OnStart hook executing: ","callee":"stopFunc","caller":"callerFunc"}`,
		},
		{
			name: "OnStopExecuted",
			event: &fxevent.OnStopExecuted{
				FunctionName: "stopFunc",
				CallerName:   "callerFunc",
				Runtime:      time.Second,
			},
			want: `{"level":"debug","message":"OnStart hook executing: ","callee":"stopFunc","caller":"callerFunc"}`,
		},
		{
			name: "OnStopExecuted with err",
			event: &fxevent.OnStopExecuted{
				FunctionName: "stopFunc",
				CallerName:   "callerFunc",
				Runtime:      time.Second,
				Err:          errors.New("error"),
			},
			want: `{"level":"debug","message":"OnStart hook failed: ","callee":"stopFunc","caller":"callerFunc","error":"error"}`,
		},
		{
			name: "Supplied",
			event: &fxevent.Supplied{
				TypeName: "typeA",
				Err:      errors.New("error"),
			},
			want: `{"level":"debug","message":"supplied: ","type":"typeA","error":"error"}`,
		},
		{
			name: "Provided",
			event: &fxevent.Provided{
				ConstructorName: "constructorA",
				OutputTypeNames: []string{"typeA", "typeB"},
			},
			want: `{"level":"debug","message":"provided: ","constructor":"constructorA","types":["typeA","typeB"]}`,
		},
		{
			name: "Decorated",
			event: &fxevent.Decorated{
				DecoratorName:   "decoratorA",
				OutputTypeNames: []string{"typeA", "typeB"},
			},
			want: `{"level":"debug","message":"decorated: ","decorator":"decoratorA","types":["typeA","typeB"]}`,
		},
		{
			name: "Invoking",
			event: &fxevent.Invoking{
				FunctionName: "invokeFunc",
			},
			want: `{"level":"debug","message":"invoking: ","callee":"invokeFunc"}`,
		},
		{
			name: "Started",
			event: &fxevent.Started{
				Err: nil,
			},
			want: `{"level":"debug","message":"Started"}`,
		},
		{
			name: "LoggerInitialized",
			event: &fxevent.LoggerInitialized{
				ConstructorName: "constructorA",
				Err:             nil,
			},
			want: `{"level":"debug","message":"Logger initialized: ","constructor":"constructorA"}`,
		},
	}
	var wg sync.WaitGroup
	wg.Add(len(tests))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				defer wg.Done()
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					if scanner.Text() == tt.want {
					} else {
						t.Errorf("LogEvent() = %v, want %v", scanner.Text(), tt.want)
					}
				}

				if scanner.Err() != nil {
					t.Errorf("LogEvent() error = %v", scanner.Err())
				}
			}()

			fxLogger.LogEvent(tt.event)
		})
	}

	wg.Wait()
}
