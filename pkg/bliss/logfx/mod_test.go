package logfx_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/logfx"
	"github.com/stretchr/testify/assert"
)

func TestRegisterLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      *logfx.Config
		wantErr     bool
		expectedErr string
	}{
		{
			name: "ValidConfig",
			config: &logfx.Config{
				Level:      "info",
				PrettyMode: true,
				AddSource:  true,
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "InvalidLogLevel",
			config: &logfx.Config{
				Level:      "invalid",
				PrettyMode: true,
				AddSource:  true,
			},
			wantErr:     true,
			expectedErr: "failed to parse log level: slog: level string \"invalid\": unknown name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logfx.RegisterLogger(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
				assert.Equal(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
			}
		})
	}
}
