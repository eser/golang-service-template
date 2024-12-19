package shared

type Cursor struct {
	Offset string `json:"offset"`
	Limit  int32  `json:"limit"`
}

type RecordsWithCursor[T any] struct {
	Cursor  string `json:"cursor"`
	Records []T    `json:"records"`
}

func NewRecordsWithCursor[T any](records []T, targetLength int32, cursorKey func(*T) string) *RecordsWithCursor[T] {
	var cursorValue string

	if len(records) == int(targetLength) { // && targetLength > 0
		cursorValue = cursorKey(&records[len(records)-1])
	}

	return &RecordsWithCursor[T]{
		Records: records,
		Cursor:  cursorValue,
	}
}
