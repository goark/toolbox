package nasaapod

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/nasaapi"
	"github.com/goark/toolbox/values"
	"github.com/goark/webinfo"
)

const (
	MediaImage = "image"
	MediaVideo = "video"
)

// Response is response data from NASA APOD API.
type Response struct {
	Date         values.Date `json:"date,omitempty"`           // APOD date in YYYY-MM-DD format.
	PostID       int         `json:"post_id,omitempty"`        // WordPress post ID.
	Title        string      `json:"title,omitempty"`          // APOD post title.
	Permalink    string      `json:"permalink,omitempty"`      // URL of the APOD post on the site.
	MediaType    string      `json:"media_type,omitempty"`     // Normalized APOD media type, such as "image", "video", or "iframe".
	Explanation  string      `json:"explanation,omitempty"`    // APOD explanation text.
	Credit       string      `json:"credit,omitempty"`         // APOD credit information.
	Copyright    string      `json:"copyright,omitempty"`      // APOD copyright information.
	Alt          string      `json:"alt,omitempty"`            // APOD alt text for the image.
	Url          string      `json:"url,omitempty"`            // APOD post URL. This matches the permalink.
	HdUrl        string      `json:"hdurl,omitempty"`          // Full-size featured image URL when available.
	BasicHTML    string      `json:"basic_html,omitempty"`     // APOD Basic HTML document as a JSON string.
	BasicHTMLUrl string      `json:"basic_html_url,omitempty"` // Plain-text HTML source URL for easier copy/paste workflows.
}

// ResponseError represents an error response from the NASA APOD API.
type ResponseError struct {
	Code    string `json:"code"`    // Error code from the API.
	Message string `json:"message"` // Error message from the API.
	Data    struct {
		Status int `json:"status"` // HTTP status code from the API.
	} `json:"data"`
}

// decode decodes APOD API responses and normalizes them into a response slice.
//
// The API may return either a single object or an array depending on endpoint
// and query mode. This function also detects structured API error payloads and
// converts them into ErrAPODAPIResponse.
func decode(r io.Reader, isSingle bool) ([]Response, error) {
	// try decode error response first
	buf := &bytes.Buffer{}
	var respErr ResponseError
	if err := json.NewDecoder(io.TeeReader(r, buf)).Decode(&respErr); err == nil {
		if isAPODErrorResponse(&respErr) {
			return nil, errs.Wrap(nasaapi.ErrAPODAPIResponse, errs.WithContext("response", respErr))
		}
	}
	// decode normal response
	var resps []Response
	dec := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	if isSingle {
		for {
			var resp Response
			if err := dec.Decode(&resp); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return nil, errs.Wrap(err)
			}
			resps = append(resps, resp)
		}
	} else {
		for {
			var resp []Response
			if err := dec.Decode(&resp); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return nil, errs.Wrap(err)
			}
			resps = append(resps, resp...)
		}
	}
	return resps, nil
}

// isAPODErrorResponse reports whether decoded JSON matches APOD API error shape.
func isAPODErrorResponse(respErr *ResponseError) bool {
	if respErr == nil {
		return false
	}
	if respErr.Data.Status >= http.StatusBadRequest {
		return true
	}
	if respErr.Data.Status == 0 && len(respErr.Code) > 0 {
		return true
	}
	return false
}

// Encode method writes encoded response data to writer by JSON format.
func (res *Response) Encode(w io.Writer) error {
	if res == nil {
		return nil
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// WebPage method returns web page of APOD.
func (res *Response) WebPage() string {
	if res == nil {
		return ""
	}
	if len(res.Permalink) > 0 {
		return res.Permalink
	}
	return res.Url
}

// ImageFile downloads the APOD image and stores it under dir.
//
// The method prefers HdUrl and falls back to Url when HdUrl is unavailable.
func (res *Response) ImageFile(ctx context.Context, dir string) (tname string, err error) {
	if res == nil {
		err = errs.Wrap(ecode.ErrNullPointer)
		return
	}
	urlStr := res.HdUrl
	if len(urlStr) == 0 {
		urlStr = res.Url
	}
	if len(urlStr) == 0 {
		err = errs.Wrap(ecode.ErrNoAPODImage)
		return
	}
	wi := &webinfo.Webinfo{ImageURL: urlStr, UserAgent: webinfo.DefaultUserAgent()}
	tname, err = wi.DownloadImage(ctx, dir, true)
	return
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
