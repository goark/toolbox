package nasaapod

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/goark/errs"
	"github.com/goark/fetch"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/values"
	"golang.org/x/net/context"
)

const (
	MediaImage = "image"
	MediaVideo = "video"
)

// Response is response data from NASA APOD API.
type Response struct {
	Copyright      string      `json:"copyright,omitempty"`
	Date           values.Date `json:"date,omitempty"`
	Explanation    string      `json:"explanation,omitempty"`
	HdUrl          string      `json:"hdurl,omitempty"`
	MediaType      string      `json:"media_type,omitempty"`
	ServiceVersion string      `json:"service_version,omitempty"`
	Title          string      `json:"title,omitempty"`
	Url            string      `json:"url,omitempty"`
	ThumbnailUrl   string      `json:"thumbnail_url,omitempty"`
}

func decode(r io.Reader, isSingle bool) ([]*Response, error) {
	var resps []*Response
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
			resps = append(resps, &resp)
		}
	} else {
		for {
			var resp []*Response
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
	u, err := url.Parse(webPage)
	if err != nil {
		return ""
	}
	u.Path = path.Join(u.Path, "ap"+res.Date.Format("060102")+".html")
	return u.String()
}

func (res *Response) ImageFile(ctx context.Context, dir string) (string, error) {
	if res == nil {
		return "", errs.Wrap(ecode.ErrNullPointer)
	}
	var urlStr string
	if res.MediaType == MediaImage {
		if len(res.Url) > 0 {
			urlStr = res.Url
		} else if len(res.HdUrl) > 0 {
			urlStr = res.HdUrl
		} else if len(res.ThumbnailUrl) > 0 {
			urlStr = res.ThumbnailUrl
		}
	} else if len(res.ThumbnailUrl) > 0 {
		urlStr = res.ThumbnailUrl
	}
	if len(urlStr) == 0 {
		return "", errs.Wrap(ecode.ErrNoAPODImage)
	}

	// get Image data
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	img, err := fetch.New().GetWithContext(ctx, u)
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	defer img.Close()

	// copy to temporary file
	file, err := os.CreateTemp(dir, "apod.*.bin")
	if err != nil {
		return "", errs.Wrap(err)
	}
	defer file.Close()

	tname := file.Name()
	_, err = io.Copy(file, img.Body())
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("url", urlStr), errs.WithContext("temp_file", tname))
	}
	return tname, nil
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
