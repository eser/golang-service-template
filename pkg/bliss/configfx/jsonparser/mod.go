package jsonparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ParseBytes(data []byte, out *map[string]any) error {
	return json.Unmarshal(data, out) //nolint:wrapcheck
}

func Parse(m *map[string]any, r io.Reader) error {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, r)
	if err != nil {
		return fmt.Errorf("parsing error: %w", err)
	}

	return ParseBytes(buf.Bytes(), m)
}

func TryParseFiles(m *map[string]any, filenames ...string) error {
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return err //nolint:wrapcheck
		}

		defer file.Close()

		err = Parse(m, file)
		if err != nil {
			return err
		}
	}

	return nil
}
