package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type MonthYear time.Time

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       *int       `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date"`
}

func (m MonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(m)
	return []byte(`"` + t.Format("01-2006") + `"`), nil
}

func (m *MonthYear) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	*m = MonthYear(t)
	return nil
}

func (m MonthYear) Time() time.Time {
	return time.Time(m)
}

func (m MonthYear) PtrTime() *time.Time {
	t := time.Time(m)
	return &t
}

func (m *MonthYear) Scan(value interface{}) error {
	if value == nil {
		*m = MonthYear(time.Time{})
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan %T into MonthYear", value)
	}
	*m = MonthYear(t)
	return nil
}

func (m MonthYear) Value() (driver.Value, error) {
	t := time.Time(m)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}
