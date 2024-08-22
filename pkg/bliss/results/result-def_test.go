package results_test

import (
	"errors"
	"testing"

	"github.com/eser/go-service/pkg/bliss/results"
	"github.com/stretchr/testify/assert"
)

var (
	resultOk  = results.NewResultDef("0001", "OK")    //nolint:gochecknoglobals
	resultErr = results.NewResultDef("0002", "Error") //nolint:gochecknoglobals
)

func TestResultSimple(t *testing.T) {
	t.Parallel()

	occurrenceOk := resultOk.New()

	assert.True(t, occurrenceOk.IsOk())
}

func TestResultError(t *testing.T) {
	t.Parallel()

	err := errors.New("error") //nolint:err113
	occurrenceErr := resultErr.NewWithError(err)

	assert.False(t, occurrenceErr.IsOk())
}
