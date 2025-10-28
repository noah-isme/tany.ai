package models

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

// JSONB wraps a map for storing arbitrary structured metadata in PostgreSQL JSONB columns.
type JSONB map[string]any

// Value implements driver.Valuer.
func (j JSONB) Value() (driver.Value, error) {
    if j == nil {
        return []byte("{}"), nil
    }
    data, err := json.Marshal(j)
    if err != nil {
        return nil, fmt.Errorf("jsonb marshal: %w", err)
    }
    return data, nil
}

// Scan implements sql.Scanner.
func (j *JSONB) Scan(value any) error {
    if j == nil {
        return fmt.Errorf("jsonb: Scan on nil pointer")
    }
    if value == nil {
        *j = JSONB{}
        return nil
    }
    var data []byte
    switch v := value.(type) {
    case string:
        data = []byte(v)
    case []byte:
        data = v
    default:
        return fmt.Errorf("jsonb: unsupported type %T", value)
    }
    if len(data) == 0 {
        *j = JSONB{}
        return nil
    }
    var raw map[string]any
    if err := json.Unmarshal(data, &raw); err != nil {
        return fmt.Errorf("jsonb unmarshal: %w", err)
    }
    *j = JSONB(raw)
    return nil
}
