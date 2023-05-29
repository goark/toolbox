package nasaapod

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"time"

	"github.com/goark/errs"
	"github.com/goark/toolbox/nasaapi"
	"github.com/goark/toolbox/values"
)

// Request is for context of APOD API.
type Request struct {
	Date      values.Date `json:"date,omitempty"`       // The date of the APOD image to retrieve
	StartDate values.Date `json:"start_date,omitempty"` // The start of a date range, when requesting date for a range of dates. Cannot be used with date.
	EndDate   values.Date `json:"end_date,omitempty"`   // The end of the date range, when used with start_date.
	Count     int         `json:"count,omitempty"`      // If this is specified then count randomly chosen images will be returned. Cannot be used with date or start_date and end_date.
	Thumbs    bool        `json:"thumbs,omitempty"`     // Return the URL of video thumbnail. If an APOD is not a video, this parameter is ignored.
	APIKey    string      `json:"api_key"`              // api.nasa.gov key for expanded usage
}

type Opts func(*Request)

// New returns new Request instance for APOD API.
func New(opts ...Opts) *Request {
	ctx := &Request{}
	for _, opt := range opts {
		opt(ctx)
	}
	return ctx
}

// WithDate returns function for setting Request.Date.
func WithDate(date values.Date) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.Date = date
		}
	}
}

// WithStartDate returns function for setting Request.StartDate.
func WithStartDate(startDate values.Date) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.StartDate = startDate
		}
	}
}

// WithEndDate returns function for setting Request.EndDate.
func WithEndDate(endDate values.Date) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.EndDate = endDate
		}
	}
}

// WithCount returns function for setting Request.Count.
func WithCount(count int) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.Count = count
		}
	}
}

// WithThumbs returns function for setting Request.Thumbs.
func WithThumbs(thumbs bool) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.Thumbs = thumbs
		}
	}
}

// WithAPIKey returns function for setting Request.APIKey.
func WithAPIKey(apiKey string) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.APIKey = apiKey
		}
	}
}

// Encode returns JSON string.
func (req *Request) Encode() (string, error) {
	if req == nil {
		return "", errs.Wrap(nasaapi.ErrNullPointer)
	}
	b, err := json.Marshal(req)
	if err != nil {
		return "", errs.Wrap(err)
	}
	return string(b), err
}

// Stringger method.
func (req *Request) String() string {
	s, err := req.Encode()
	if err != nil {
		return ""
	}
	return s
}

// Get method gets APOD data from NASA API, and returns []*Response instance.
func (req *Request) Get(ctx context.Context) ([]*Response, error) {
	if req == nil {
		return nil, errs.Wrap(nasaapi.ErrNullPointer)
	}
	resp, err := req.getRawData(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	return decode(resp, req.isSingle())
}

func (req *Request) getRawData(ctx context.Context) (io.ReadCloser, error) {
	if req == nil {
		return nil, errs.Wrap(nasaapi.ErrNullPointer)
	}
	q, err := req.makeQuery()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return nasaapi.Fetch(ctx, apiPath, q)
}

func (req *Request) isSingle() bool {
	if !req.Date.IsZero() {
		return true
	}
	if req.StartDate.IsZero() && req.EndDate.IsZero() && req.Count == 0 {
		return true
	}
	return false
}

func (req *Request) makeQuery() (url.Values, error) {
	v := url.Values{}
	if !req.Date.IsZero() {
		if !req.StartDate.IsZero() || !req.EndDate.IsZero() || req.Count > 0 {
			return nil, errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
		v.Set("date", req.Date.Format(time.DateOnly))
	}
	if !req.StartDate.IsZero() {
		if !req.Date.IsZero() || req.Count > 0 {
			return nil, errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
		v.Set("start_date", req.StartDate.Format(time.DateOnly))
	}
	if !req.EndDate.IsZero() {
		if req.StartDate.IsZero() || !req.Date.IsZero() || req.Count > 0 {
			return nil, errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
		v.Set("end_date", req.EndDate.Format(time.DateOnly))
	}
	if req.Count > 0 {
		v.Set("count", strconv.Itoa(req.Count))
	}
	if req.Thumbs {
		v.Set("thumbs", "true")
	}
	if len(req.APIKey) > 0 {
		v.Set("api_key", req.APIKey)
	} else {
		v.Set("api_key", nasaapi.DefaultAPIKey)
	}
	return v, nil
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
