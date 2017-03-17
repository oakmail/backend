package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Resource is a basic storage in the system without any particular structure.
type Resource struct {
	ID           uint64    `db:"id" json:"id" goqu:"skipinsert"`
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	DateModified time.Time `db:"date_modified" json:"date_modified"`
	Owner        uint64    `db:"owner" json:"owner"`
	Meta         Meta      `db:"meta" json:"meta"`
	Tags         Tags      `db:"tags" json:"tags"`
	File         string    `db:"file" json:"file,omitempty"`
	UploadToken  string    `db:"upload_token" json:"upload_token,omitempty"`
}

// Meta is a custom JSON storage type for objects
type Meta map[string]interface{}

// Scan implements the database/sql Scanner interface
func (m *Meta) Scan(value interface{}) error {
	var data string
	if err := json.Unmarshal(append(append([]byte{'"'}, value.([]byte)...), '"'), &data); err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), m)
}

// Value implements the database/sql Valuer interface
func (m Meta) Value() (driver.Value, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(data), err
}

// Tags is a custom JSON storage type for arrays
type Tags []string

// Scan implements the database/sql Scanner interface
func (t *Tags) Scan(value interface{}) error {
	var data string
	if err := json.Unmarshal(append(append([]byte{'"'}, value.([]byte)...), '"'), &data); err != nil {
		return err
	}

	if data[0] != '[' || data[len(data)-1] != ']' {
		data = `["` + data + `"]`
	}

	return json.Unmarshal([]byte(data), t)
}

// Value implements the database/sql Valuer interface
func (t Tags) Value() (driver.Value, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return string(data), err
}
