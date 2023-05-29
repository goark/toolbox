package values

import (
	"strconv"
	"strings"
	"time"

	"github.com/goark/errs"
)

// Date is wrapper class of time.Time.
type Date struct {
	time.Time
}

// NewDate returns Date instance from time.Time.
func NewDate(tm time.Time) Date {
	return Date{tm}
}

// Today returns Date instance in today.
func Today() Date {
	return NewDate(time.Now())
}

// Stringer with YYYY-MM-DD format.
func (t Date) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateOnly)
}

// MarshalJSON implements the json.Marshaler interface.
func (t Date) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(strconv.Quote(t.String())), nil
}

var timeTemplate = []string{
	time.RFC3339,
	time.DateOnly,
}

// DateFrom returns Date instance from date string
func DateFrom(s string) (Date, error) {
	if len(s) == 0 || strings.EqualFold(s, "null") {
		return NewDate(time.Time{}), nil
	}
	var lastErr error
	for _, tmplt := range timeTemplate {
		if tm, err := time.Parse(tmplt, s); err != nil {
			lastErr = errs.Wrap(err, errs.WithContext("time_string", s), errs.WithContext("time_template", tmplt))
		} else {
			return NewDate(tm), nil
		}
	}
	return NewDate(time.Time{}), lastErr
}

// UnmarshalJSON implements the json.UnmarshalJSON interface.
func (t *Date) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		s = string(b)
	}
	*t, err = DateFrom(s)
	return err
}

// Equal return true if left equals right in year/month/day.
func (left Date) Equal(right Date) bool {
	if left.IsZero() || right.IsZero() {
		return false
	}
	r := right.In(left.Location())
	return left.Year() == r.Year() && left.YearDay() == r.YearDay()
}

// After return true if left is not equal right and left > right.
func (left Date) After(right Date) bool {
	if left.IsZero() || right.IsZero() {
		return false
	}
	if left.Equal(right) {
		return false
	}
	r := right.In(left.Location())
	return left.Time.After(r)
}

// Aftere return true if left is not equal right and left < right.
func (left Date) Before(right Date) bool {
	if left.IsZero() || right.IsZero() {
		return false
	}
	if left.Equal(right) {
		return false
	}
	return !left.After(right)
}

/* Copyright 2023 Spiegel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
