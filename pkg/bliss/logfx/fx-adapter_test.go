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
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx/fxevent"
)

type MockWriter struct{}

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

	assert.NotNil(t, fxLogger, "GetFxLogger() = nil, want not nil")
}

func TestFxLogger_LogEvent(t *testing.T) { //nolint:paralleltest
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

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				defer wg.Done()

				scanner := bufio.NewScanner(os.Stdin)

				for scanner.Scan() {
					assert.Equal(t, tt.want, scanner.Text())
				}

				assert.NoError(t, scanner.Err())
			}()

			fxLogger.LogEvent(tt.event)
		})
	}

	wg.Wait()
}
