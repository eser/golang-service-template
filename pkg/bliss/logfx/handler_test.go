package logfx_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/eser/go-service/pkg/bliss/logfx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFailWriter struct{}

func (m *mockFailWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to write")
}

func TestNewHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		writer      *bytes.Buffer
		config      *logfx.Config
		expectedErr string
	}{
		{
			name:   "ValidConfig",
			writer: &bytes.Buffer{},
			config: &logfx.Config{
				Level:      "info",
				PrettyMode: true,
				AddSource:  true,
			},
			expectedErr: "",
		},
		{
			name:   "InvalidLogLevel",
			writer: &bytes.Buffer{},
			config: &logfx.Config{
				Level:      "invalid",
				PrettyMode: true,
				AddSource:  true,
			},
			expectedErr: "failed to parse log level: slog: level string \"invalid\": unknown name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler, err := logfx.NewHandler(tt.writer, tt.config)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				assert.Nil(t, handler)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, handler)
		})
	}
}

func TestHandler_Enabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		level    string
		expected bool
	}{
		{
			name:     "Enabled",
			level:    "info",
			expected: true,
		},
		{
			name:     "Disabled",
			level:    "error",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler, _ := logfx.NewHandler(&bytes.Buffer{}, &logfx.Config{
				Level: tt.level,
			})

			assert.Equal(t, tt.expected, handler.Enabled(context.Background(), 0))
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		level    string
		record   slog.Record
		expected string
	}{
		{
			name:     "Debug",
			level:    "debug",
			record:   slog.NewRecord(time.Time{}, slog.LevelDebug, "test", 0),
			expected: "\x1b[90m00:00:00.000\x1b[0m \x1b[94mDEBUG\x1b[0m test {}\n",
		},
		{
			name:     "Info",
			level:    "info",
			record:   slog.NewRecord(time.Time{}, slog.LevelInfo, "test", 0),
			expected: "\x1b[90m00:00:00.000\x1b[0m \x1b[32mINFO\x1b[0m test {}\n",
		},
		{
			name:     "Warn",
			level:    "warn",
			record:   slog.NewRecord(time.Time{}, slog.LevelWarn, "test", 0),
			expected: "\x1b[90m00:00:00.000\x1b[0m \x1b[33mWARN\x1b[0m test {}\n",
		},
		{
			name:     "Error",
			level:    "error",
			record:   slog.NewRecord(time.Time{}, slog.LevelError, "test", 0),
			expected: "\x1b[90m00:00:00.000\x1b[0m \x1b[31mERROR\x1b[0m test {}\n",
		},
		{
			name:     "UnknownLevel",
			level:    "error",
			record:   slog.NewRecord(time.Time{}, 77, "test", 0),
			expected: "\x1b[90m00:00:00.000\x1b[0m ERROR+69 test {}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			writer := &bytes.Buffer{}
			handler, _ := logfx.NewHandler(writer, &logfx.Config{
				Level:      tt.level,
				PrettyMode: true,
			})

			err := handler.Handle(context.Background(), tt.record)
			require.NoError(t, err)

			assert.Contains(t, writer.String(), tt.expected)
		})
	}

	t.Run("failed to write log", func(t *testing.T) {
		t.Parallel()
		handler, _ := logfx.NewHandler(&mockFailWriter{}, &logfx.Config{
			Level:      "info",
			PrettyMode: true,
		})
		err := handler.Handle(context.Background(), slog.NewRecord(time.Time{}, slog.LevelInfo, "test", 0))
		assert.EqualError(t, err, "failed to write log: failed to write")
	})
}

func TestHandler_WithAttrs(t *testing.T) {
	t.Parallel()
	handler, _ := logfx.NewHandler(&bytes.Buffer{}, &logfx.Config{
		Level: "info",
	})
	newHandler := handler.WithAttrs([]slog.Attr{})
	// FIXME(@eser) should equal or not?
	assert.Equal(t, handler, newHandler)
}

func TestHandler_WithGroup(t *testing.T) {
	t.Parallel()
	handler, _ := logfx.NewHandler(&bytes.Buffer{}, &logfx.Config{
		Level: "info",
	})
	newHandler := handler.WithGroup("test")
	assert.NotEqual(t, handler, newHandler)
}
