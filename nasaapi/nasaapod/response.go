package nasaapod

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/values"
	"github.com/goark/webinfo"
)

const (
	MediaImage = "image"
	MediaVideo = "video"
)

// Response is response data from NASA APOD API.
type Response struct {
	Date        values.Date `json:"date,omitempty"`
	PostID      int         `json:"post_id,omitempty"`
	Title       string      `json:"title,omitempty"`
	Permalink   string      `json:"permalink,omitempty"`
	MediaType   string      `json:"media_type,omitempty"`
	Explanation string      `json:"explanation,omitempty"`
	Credit      string      `json:"credit,omitempty"`
	Copyright   string      `json:"copyright,omitempty"`
	Alt         string      `json:"alt,omitempty"`
	Url         string      `json:"url,omitempty"`
	HdUrl       string      `json:"hdurl,omitempty"`

	ServiceVersion string `json:"service_version,omitempty"`
}

func decode(r io.Reader, isSingle bool) ([]Response, error) {
	var resps []Response
	dec := json.NewDecoder(r)
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
