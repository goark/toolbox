package nasaapod

import (
	"errors"
	"testing"

	"github.com/goark/toolbox/nasaapi"
	"github.com/goark/toolbox/values"
)

func dateFromMust(s string) values.Date {
	dt, err := values.DateFrom(s)
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
		count     int
		thumbs    bool
		apiKey    string
		err       error
		want      string
	}{
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			count:     0,
			thumbs:    false,
			apiKey:    "",
			err:       nil,
			want:      `{"date":"","start_date":"","end_date":"","api_key":""}`,
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			count:     0,
			thumbs:    false,
			apiKey:    "",
			err:       nil,
			want:      `{"date":"2023-02-22","start_date":"","end_date":"","api_key":""}`,
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust(""),
			count:     0,
			thumbs:    false,
			apiKey:    "",
			err:       nil,
			want:      `{"date":"","start_date":"2023-02-22","end_date":"","api_key":""}`,
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust("2023-02-22"),
			count:     0,
			thumbs:    false,
			apiKey:    "",
			err:       nil,
			want:      `{"date":"","start_date":"2023-02-22","end_date":"2023-02-22","api_key":""}`,
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			count:     1,
			thumbs:    false,
			apiKey:    "",
			err:       nil,
			want:      `{"date":"","start_date":"","end_date":"","count":1,"api_key":""}`,
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			count:     0,
			thumbs:    true,
			apiKey:    "foo",
			err:       nil,
			want:      `{"date":"","start_date":"","end_date":"","thumbs":true,"api_key":"foo"}`,
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust(""),
			count:     0,
			thumbs:    false,
			apiKey:    "",
			err:       nasaapi.ErrCombination,
			want:      "",
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust(""),
			endDate:   dateFromMust("2023-02-22"),
			count:     0,
			thumbs:    false,
			apiKey:    "",
			err:       nasaapi.ErrCombination,
			want:      "",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust(""),
			endDate:   dateFromMust("2023-02-22"),
			count:     0,
			thumbs:    false,
			apiKey:    "",
			err:       nasaapi.ErrCombination,
			want:      "",
		},
		{
			date:      dateFromMust("2023-02-22"),
			startDate: dateFromMust(""),
			endDate:   dateFromMust(""),
			count:     1,
			thumbs:    false,
			apiKey:    "",
			err:       nasaapi.ErrCombination,
			want:      "",
		},
		{
			date:      dateFromMust(""),
			startDate: dateFromMust("2023-02-22"),
			endDate:   dateFromMust(""),
			count:     1,
			thumbs:    false,
			apiKey:    "",
			err:       nasaapi.ErrCombination,
			want:      "",
		},
	}

	for _, tc := range testCases {
		req := New(
			WithDate(tc.date),
			WithStartDate(tc.startDate),
			WithEndDate(tc.endDate),
			WithCount(tc.count),
			WithThumbs(tc.thumbs),
			WithAPIKey(tc.apiKey),
		)
		_, err := req.makeQuery()
		if !errors.Is(err, tc.err) {
			t.Errorf("makeQuery() is \"%v\", want \"%v\"", err, tc.err)
		}
		if err == nil {
			if got, err := req.Encode(); err != nil {
				t.Errorf("Encode() is \"%v\", want nil", err)
			} else if got != tc.want {
				t.Errorf("Encode() = \"%v\", want \"%v\"", got, tc.want)
			}

		}
	}
}

/* MIT License
 *
 * Copyright 2023 Spiegel
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
