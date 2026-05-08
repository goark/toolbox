package nasaapod

import (
	"context"
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

func (res *Response) ImageFile(ctx context.Context, dir string) (tname string, err error) {
	if res == nil {
		err = errs.Wrap(ecode.ErrNullPointer)
		return
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
		err = errs.Wrap(ecode.ErrNoAPODImage)
		return
	}

	// get Image data
	u, perr := url.Parse(urlStr)
	if perr != nil {
		err = errs.Wrap(perr, errs.WithContext("url", urlStr))
		return
	}
	img, ferr := fetch.New().GetWithContext(ctx, u)
	if ferr != nil {
		err = errs.Wrap(ferr, errs.WithContext("url", urlStr))
		return
	}
	defer func() {
		err = errs.Join(err, img.Close())
	}()

	// copy to temporary file
	file, ferr := os.CreateTemp(dir, "apod.*.bin")
	if ferr != nil {
		err = errs.Wrap(ferr)
		return
	}
	defer func() {
		err = errs.Join(err, file.Close())
	}()

	tname = file.Name()
	if _, cerr := io.Copy(file, img.Body()); cerr != nil {
		err = errs.Wrap(cerr, errs.WithContext("url", urlStr), errs.WithContext("temp_file", tname))
		return
	}
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
