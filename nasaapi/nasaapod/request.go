package nasaapod

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strconv"

	"github.com/goark/errs"
	"github.com/goark/toolbox/nasaapi"
	"github.com/goark/toolbox/values"
)

const (
	defaultPerPage = 25
	maxPerPage     = 100
	maxPage        = 500
)

// Request is for context of APOD API.
type Request struct {
	Date      values.Date `json:"date,omitempty"`       // The date of the APOD image to retrieve
	StartDate values.Date `json:"start_date,omitempty"` // The start of a date range, when requesting date for a range of dates. Cannot be used with date.
	EndDate   values.Date `json:"end_date,omitempty"`   // The end of the date range, when used with start_date.
	Count     int         `json:"count,omitempty"`      // If this is specified then count randomly chosen images will be returned. Cannot be used with date or start_date and end_date.
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

// Get method gets APOD data from APOD API.
func (req *Request) Get(ctx context.Context) (rsp []Response, err error) {
	if req == nil {
		err = errs.Wrap(nasaapi.ErrNullPointer)
		return
	}
	if err = req.validate(); err != nil {
		return nil, err
	}
	perPage := req.perPage()
	for page := 1; page <= maxPage; page++ {
		resps, ferr := req.fetchPage(ctx, page, perPage)
		if ferr != nil {
			return nil, errs.Wrap(ferr)
		}
		if len(resps) == 0 {
			break
		}
		matched := req.filterResponses(resps)
		rsp = append(rsp, matched...)
		if req.done(rsp, resps) {
			break
		}
	}
	if req.Count > 0 && len(rsp) > req.Count {
		rsp = rsp[:req.Count]
	}
	return
}

func (req *Request) fetchPage(ctx context.Context, page, perPage int) ([]Response, error) {
	resp, err := req.getRawData(ctx, page, perPage)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer func() {
		_ = resp.Close()
	}()
	resps, err := decode(resp, false)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return resps, nil
}

func (req *Request) getRawData(ctx context.Context, page, perPage int) (io.ReadCloser, error) {
	if req == nil {
		return nil, errs.Wrap(nasaapi.ErrNullPointer)
	}
	q, err := req.makeQuery(page, perPage)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return nasaapi.Fetch(ctx, nasaapi.APODHost, apiPath, q)
}

func (req *Request) validate() error {
	if req == nil {
		return errs.Wrap(nasaapi.ErrNullPointer)
	}
	if !req.Date.IsZero() {
		if !req.StartDate.IsZero() || !req.EndDate.IsZero() || req.Count > 0 {
			return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
	}
	if !req.StartDate.IsZero() {
		if !req.Date.IsZero() || req.Count > 0 {
			return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
	}
	if !req.EndDate.IsZero() {
		if req.StartDate.IsZero() || !req.Date.IsZero() || req.Count > 0 {
			return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
	}
	if req.Count < 0 {
		return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
	}
	return nil
}

func (req *Request) perPage() int {
	if req == nil || req.Count <= 0 {
		return defaultPerPage
	}
	if req.Count > maxPerPage {
		return maxPerPage
	}
	return req.Count
}

func (req *Request) makeQuery(page, perPage int) (url.Values, error) {
	if err := req.validate(); err != nil {
		return nil, errs.Wrap(err)
	}
	if page <= 0 || perPage <= 0 {
		return nil, errs.Wrap(nasaapi.ErrCombination, errs.WithContext("page", page), errs.WithContext("per_page", perPage))
	}
	v := url.Values{}
	v.Set("page", strconv.Itoa(page))
	v.Set("per_page", strconv.Itoa(perPage))
	return v, nil
}

func (req *Request) filterResponses(resps []Response) []Response {
	if req == nil {
		return nil
	}
	matched := make([]Response, 0, len(resps))
	for _, r := range resps {
		if !req.Date.IsZero() {
			if r.Date.Equal(req.Date) {
				matched = append(matched, r)
			}
			continue
		}
		if !req.StartDate.IsZero() {
			if r.Date.Before(req.StartDate) {
				continue
			}
			if !req.EndDate.IsZero() && r.Date.After(req.EndDate) {
				continue
			}
			matched = append(matched, r)
			continue
		}
		matched = append(matched, r)
	}
	return matched
}

func (req *Request) done(rsp, pageData []Response) bool {
	if req == nil {
		return true
	}
	if req.Count > 0 && len(rsp) >= req.Count {
		return true
	}
	if req.Date.IsZero() && req.StartDate.IsZero() {
		return false
	}
	if len(pageData) == 0 {
		return true
	}
	oldest := pageData[len(pageData)-1].Date
	if !req.Date.IsZero() {
		return len(rsp) > 0 || oldest.Before(req.Date)
	}
	return oldest.Before(req.StartDate)
}

/* Copyright 2023-2026 Spiegel
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
