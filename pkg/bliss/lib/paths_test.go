package lib_test

import (
	"testing"

	"github.com/eser/go-service/pkg/bliss/lib"
	"github.com/stretchr/testify/assert"
)

func TestPathsSplit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		wantDir  string
		wantFile string
		wantExt  string
	}{
		{
			name:     "Empty filename",
			filename: "",
			wantDir:  "",
			wantFile: "",
			wantExt:  "",
		},
		{
			name:     "Filename without extension",
			filename: "path/to/file",
			wantDir:  "path/to/",
			wantFile: "file",
			wantExt:  "",
		},
		{
			name:     "Filename with extension",
			filename: "path/to/file.txt",
			wantDir:  "path/to/",
			wantFile: "file",
			wantExt:  ".txt",
		},
		{
			name:     "Filename with multiple dots",
			filename: "path/to/file.with.multiple.dots.txt",
			wantDir:  "path/to/",
			wantFile: "file.with.multiple.dots",
			wantExt:  ".txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotDir, gotFile, gotExt := lib.PathsSplit(tt.filename)

			assert.Equal(t, gotDir, tt.wantDir)
			assert.Equal(t, gotFile, tt.wantFile)
			assert.Equal(t, gotExt, tt.wantExt)
		})
	}
}
