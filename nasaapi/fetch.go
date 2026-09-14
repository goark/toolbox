package nasaapi

import (
	"context"
	"io"
	"net/url"

	"github.com/goark/errs"
	"github.com/goark/fetch"
)

const (
	DefaultScheme = "https"            // Default scheme for NASA API requests
	DefaultHost   = "api.nasa.gov"     // Default host for NASA API requests
	APODHost      = "science.nasa.gov" // Host for APOD API requests
)

// Request function requests to NASA API, and returns response data.
func Fetch(ctx context.Context, host string, path string, q url.Values) (io.ReadCloser, error) {
	u := getURL(host, path, q)
	resp, err := fetch.New().GetWithContext(ctx, u)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return resp.Body(), nil
}

func getURL(host string, path string, q url.Values) *url.URL {
	return &url.URL{
		Scheme:   DefaultScheme,
		Host:     host,
		Path:     path,
		RawQuery: q.Encode(),
	}
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
