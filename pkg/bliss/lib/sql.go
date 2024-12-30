package lib

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type NullString struct { //nolint:recvcheck
	String string
	Valid  bool
}

func (ns *NullString) Scan(value any) error {
	if value == nil {
		ns.String = ""
		ns.Valid = false

		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("NullString: unexpected type %T", value) //nolint:err113
	}

	ns.String = str
	ns.Valid = true

	return nil
}

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil //nolint:nilnil
	}

	return ns.String, nil
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.String = ""
		ns.Valid = false

		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err //nolint:wrapcheck
	}

	ns.String = s
	ns.Valid = true

	return nil
}
