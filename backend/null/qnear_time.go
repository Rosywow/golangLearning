package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// QNearTime is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type QNearTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (t *QNearTime) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	t.Valid = (err == nil)
	return err
}

// Value implements the driver Valuer interface.
func (t QNearTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewQNearTime creates a new Time.
func NewQNearTime(t time.Time, valid bool) QNearTime {
	return QNearTime{
		Time:  t,
		Valid: valid,
	}
}

// QNearTimeFrom creates a new Time that will always be valid.
func QNearTimeFrom(t time.Time) QNearTime {
	return NewQNearTime(t, true)
}

// QNearTimeFromPtr creates a new Time that will be null if t is nil.
func QNearTimeFromPtr(t *time.Time) QNearTime {
	if t == nil {
		return NewQNearTime(time.Time{}, false)
	}
	return NewQNearTime(*t, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t QNearTime) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t QNearTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}

	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *QNearTime) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		err = t.Time.UnmarshalJSON(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

//MarshalText customize version
func (t QNearTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

//UnmarshalText comstomize version
func (t *QNearTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// SetValid changes this Time's value and sets it to be non-null.
func (t *QNearTime) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t QNearTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for invalid Times, hopefully for future omitempty support.
// A non-null Time with a zero value will not be considered zero.
func (t QNearTime) IsZero() bool {
	return !t.Valid
}
