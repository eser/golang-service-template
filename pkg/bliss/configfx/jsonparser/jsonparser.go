package jsonparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ParseBytes(data []byte, out *map[string]any) error {
	var raw map[string]any

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("parsing error: %w", err)
	}

	flattenJSON(raw, "", out)

	return nil
}

func Parse(m *map[string]any, r io.Reader) error {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, r)
	if err != nil {
		return fmt.Errorf("parsing error: %w", err)
	}

	return ParseBytes(buf.Bytes(), m)
}

func tryParseFile(m *map[string]any, filename string) (err error) {
	file, fileErr := os.Open(filepath.Clean(filename))
	if fileErr != nil {
		if os.IsNotExist(fileErr) {
			return nil
		}

		return fmt.Errorf("parsing error: %w", fileErr)
	}

	defer func() {
		err = file.Close()
	}()

	return Parse(m, file)
}

func TryParseFiles(m *map[string]any, filenames ...string) error {
	for _, filename := range filenames {
		err := tryParseFile(m, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func flattenJSON(input map[string]any, prefix string, out *map[string]any) {
	for key, value := range input {
		mapValue, isMap := value.(map[string]any)

		if isMap {
			// Eğer değer bir map ise, recursive olarak çağırıyoruz
			flattenJSON(mapValue, prefix+strings.ToUpper(key)+"__", out)

			continue
		}

		// Eğer değer map değilse, anahtarı ekliyoruz
		(*out)[prefix+strings.ToUpper(key)] = fmt.Sprintf("%v", value)
	}
}
