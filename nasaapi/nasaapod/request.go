package nasaapod

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"path"
	"strconv"

	"github.com/goark/errs"
	"github.com/goark/toolbox/nasaapi"
	"github.com/goark/toolbox/values"
)

const (
	defaultPerPage = 25
	maxPerPage     = defaultPerPage
)

// Request is for context of APOD API.
type Request struct {
	Date      values.Date `json:"date,omitempty"`       // Single APOD date (YYYY-MM-DD), resolved to YYMMDD path.
	StartDate values.Date `json:"start_date,omitempty"` // Date range start (YYYY-MM-DD), encoded as date_from=YYMMDD.
	EndDate   values.Date `json:"end_date,omitempty"`   // Date range end (YYYY-MM-DD), encoded as date_to=YYMMDD.
	Page      int         `json:"page,omitempty"`       // Page number for list results.
	PerPage   int         `json:"per_page,omitempty"`   // Number of results per page (max 25).
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

// WithPage returns function for setting Request.Page.
func WithPage(page int) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.Page = page
		}
	}
}

// WithPerPage returns function for setting Request.PerPage.
func WithPerPage(perPage int) Opts {
	return func(ctx *Request) {
		if ctx != nil {
			ctx.PerPage = perPage
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
// Deprecated: Prefer using GetList or GetByDate for explicit request mode.
func (req *Request) Get(ctx context.Context) (rsp []Response, err error) {
	if req == nil {
		err = errs.Wrap(nasaapi.ErrNullPointer)
		return
	}
	if req.Date.IsZero() {
		return req.GetList(ctx)
	}
	resp, err := req.GetByDate(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if resp == nil {
		return nil, nil
	}
	return []Response{*resp}, nil
}

// GetList gets APOD list data from APOD API.
func (req *Request) GetList(ctx context.Context) (rsp []Response, err error) {
	if req == nil {
		err = errs.Wrap(nasaapi.ErrNullPointer)
		return
	}
	if err = req.validateListMode(); err != nil {
		return nil, err
	}
	resps, ferr := req.fetchList(ctx)
	if ferr != nil {
		return nil, errs.Wrap(ferr)
	}
	return resps, nil
}

// GetByDate gets APOD single item data from APOD API by Request.Date.
func (req *Request) GetByDate(ctx context.Context) (*Response, error) {
	if req == nil {
		return nil, errs.Wrap(nasaapi.ErrNullPointer)
	}
	if err := req.validateDateMode(); err != nil {
		return nil, errs.Wrap(err)
	}
	resp, err := req.fetchByDate(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return resp, nil
}

func (req *Request) fetchByDate(ctx context.Context) (*Response, error) {
	resp, err := req.getRawDataByDate(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer func() {
		_ = resp.Close()
	}()
	resps, err := decode(resp, true)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if len(resps) == 0 {
		return nil, nil
	}
	return &resps[0], nil
}

func (req *Request) fetchList(ctx context.Context) ([]Response, error) {
	resp, err := req.getRawDataList(ctx)
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

func (req *Request) getRawDataByDate(ctx context.Context) (io.ReadCloser, error) {
	if req == nil {
		return nil, errs.Wrap(nasaapi.ErrNullPointer)
	}
	return nasaapi.Fetch(ctx, nasaapi.APODHost, path.Join(apiPath, req.dateCode(req.Date)), nil)
}

func (req *Request) getRawDataList(ctx context.Context) (io.ReadCloser, error) {
	if req == nil {
		return nil, errs.Wrap(nasaapi.ErrNullPointer)
	}
	q, err := req.makeQuery()
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
		if !req.StartDate.IsZero() || !req.EndDate.IsZero() || req.Page > 0 || req.PerPage > 0 {
			return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
	}
	if !req.StartDate.IsZero() {
		if !req.Date.IsZero() {
			return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
	}
	if !req.EndDate.IsZero() {
		if req.StartDate.IsZero() || !req.Date.IsZero() {
			return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
		}
	}
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() && req.StartDate.After(req.EndDate) {
		return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
	}
	if req.Page < 0 || req.PerPage < 0 {
		return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
	}
	return nil
}

func (req *Request) validateDateMode() error {
	if req == nil {
		return errs.Wrap(nasaapi.ErrNullPointer)
	}
	if req.Date.IsZero() {
		return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
	}
	if !req.StartDate.IsZero() || !req.EndDate.IsZero() || req.Page > 0 || req.PerPage > 0 {
		return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
	}
	return req.validate()
}

func (req *Request) validateListMode() error {
	if req == nil {
		return errs.Wrap(nasaapi.ErrNullPointer)
	}
	if !req.Date.IsZero() {
		return errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
	}
	return req.validate()
}

func (req *Request) page() int {
	if req == nil || req.Page <= 0 {
		return 1
	}
	return req.Page
}

func (req *Request) perPage() int {
	if req == nil || req.PerPage <= 0 {
		return defaultPerPage
	}
	if req.PerPage > maxPerPage {
		return maxPerPage
	}
	return req.PerPage
}

func (req *Request) makeQuery() (url.Values, error) {
	if err := req.validate(); err != nil {
		return nil, errs.Wrap(err)
	}
	if !req.Date.IsZero() {
		return nil, errs.Wrap(nasaapi.ErrCombination, errs.WithContext("config", req))
	}
	v := url.Values{}
	v.Set("page", strconv.Itoa(req.page()))
	v.Set("per_page", strconv.Itoa(req.perPage()))
	if !req.StartDate.IsZero() {
		v.Set("date_from", req.dateCode(req.StartDate))
	}
	if !req.EndDate.IsZero() {
		v.Set("date_to", req.dateCode(req.EndDate))
	}
	return v, nil
}

func (req *Request) dateCode(date values.Date) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("060102")
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
