package nasaapod

import (
	"errors"
	"testing"

	"github.com/goark/toolbox/nasaapi"
	"github.com/goark/toolbox/values"
)

func dateFromMust(s string) values.Date {
	dt, err := values.DateFrom(s, false)
	if err != nil {
		panic(err)
	}
	return dt
}

func TestDate(t *testing.T) {
	testCases := []struct {
		date      values.Date
		startDate values.Date
		endDate   values.Date
		page      int
		perPage   int
		err       error
		want      string
		wantQuery string
	}{
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			page:      0,
			perPage:   0,
			err:       nil,
			want:      `{"date":"","start_date":"","end_date":""}`,
			wantQuery: "page=1&per_page=25",
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			page:      0,
			perPage:   0,
			err:       nil,
			want:      `{"date":"2023-02-22","start_date":"","end_date":""}`,
			wantQuery: "",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust(""),
			page:      0,
			perPage:   0,
			err:       nil,
			want:      `{"date":"","start_date":"2023-02-22","end_date":""}`,
			wantQuery: "date_from=230222&page=1&per_page=25",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust("2023-02-22"),
			page:      0,
			perPage:   0,
			err:       nil,
			want:      `{"date":"","start_date":"2023-02-22","end_date":"2023-02-22"}`,
			wantQuery: "date_from=230222&date_to=230222&page=1&per_page=25",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			page:      2,
			perPage:   10,
			err:       nil,
			want:      `{"date":"","start_date":"","end_date":"","page":2,"per_page":10}`,
			wantQuery: "page=2&per_page=10",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			page:      0,
			perPage:   40,
			err:       nil,
			want:      `{"date":"","start_date":"","end_date":"","per_page":40}`,
			wantQuery: "page=1&per_page=25",
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust(""),
			page:      0,
			perPage:   0,
			err:       nasaapi.ErrCombination,
			want:      "",
			wantQuery: "",
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust(""),
			endDate:   dateFromMust("2023-02-22"),
			page:      0,
			perPage:   0,
			err:       nasaapi.ErrCombination,
			want:      "",
			wantQuery: "",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust("2023-02-22"),
			page:      0,
			perPage:   0,
			err:       nasaapi.ErrCombination,
			want:      "",
			wantQuery: "",
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			page:      1,
			perPage:   0,
			err:       nasaapi.ErrCombination,
			want:      "",
			wantQuery: "",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust("2023-02-21"),
			page:      0,
			perPage:   0,
			err:       nasaapi.ErrCombination,
			want:      "",
			wantQuery: "",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			page:      -1,
			perPage:   0,
			err:       nasaapi.ErrCombination,
			want:      "",
			wantQuery: "",
		},
	}

	for _, tc := range testCases {
		req := New(
			WithDate(tc.date),
			WithStartDate(tc.startDate),
			WithEndDate(tc.endDate),
			WithPage(tc.page),
			WithPerPage(tc.perPage),
		)
		if !tc.date.IsZero() {
			err := req.validate()
			if !errors.Is(err, tc.err) {
				t.Errorf("validate() is \"%v\", want \"%v\"", err, tc.err)
			}
			if err == nil {
				if got, err := req.Encode(); err != nil {
					t.Errorf("Encode() is \"%v\", want nil", err)
				} else if got != tc.want {
					t.Errorf("Encode() = \"%v\", want \"%v\"", got, tc.want)
				}
			}
			continue
		}
		q, err := req.makeQuery()
		if !errors.Is(err, tc.err) {
			t.Errorf("makeQuery() is \"%v\", want \"%v\"", err, tc.err)
		}
		if err == nil {
			if got := q.Encode(); got != tc.wantQuery {
				t.Errorf("query = %q, want %q", got, tc.wantQuery)
			}
			if got, err := req.Encode(); err != nil {
				t.Errorf("Encode() is \"%v\", want nil", err)
			} else if got != tc.want {
				t.Errorf("Encode() = \"%v\", want \"%v\"", got, tc.want)
			}

		}
	}
}

func TestValidateDateMode(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
		err  error
	}{
		{
			name: "ok",
			req: New(
				WithDate(dateFromMust("2026-09-19")),
			),
			err: nil,
		},
		{
			name: "no date",
			req:  New(),
			err:  nasaapi.ErrCombination,
		},
		{
			name: "date with page",
			req: New(
				WithDate(dateFromMust("2026-09-19")),
				WithPage(2),
			),
			err: nasaapi.ErrCombination,
		},
		{
			name: "date with range",
			req: New(
				WithDate(dateFromMust("2026-09-19")),
				WithStartDate(dateFromMust("2026-09-01")),
			),
			err: nasaapi.ErrCombination,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.validateDateMode()
			if !errors.Is(err, tc.err) {
				t.Errorf("validateDateMode() is %v, want %v", err, tc.err)
			}
		})
	}
}

func TestValidateListMode(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
		err  error
	}{
		{
			name: "ok default",
			req:  New(),
			err:  nil,
		},
		{
			name: "ok range",
			req: New(
				WithStartDate(dateFromMust("2026-09-01")),
				WithEndDate(dateFromMust("2026-09-19")),
			),
			err: nil,
		},
		{
			name: "contains date",
			req: New(
				WithDate(dateFromMust("2026-09-19")),
			),
			err: nasaapi.ErrCombination,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.validateListMode()
			if !errors.Is(err, tc.err) {
				t.Errorf("validateListMode() is %v, want %v", err, tc.err)
			}
		})
	}
}

/* MIT License
 *
 * Copyright 2023-2026 Spiegel
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
