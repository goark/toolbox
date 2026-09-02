package webpage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/webinfo"
)

// Webpage is information of web page
type Webpage struct {
	URL         string     `json:"url,omitempty"`
	Canonical   string     `json:"canonical,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	ImageURL    string     `json:"image_url,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
}

// ReadPage function reads web page from URL, and analysis information.
func ReadPage(ctx context.Context, urlStr string) (link *Webpage, err error) {
	wi, ferr := webinfo.Fetch(ctx, urlStr, "") //default user-agent is set in webinfo.Fetch() function
	if ferr != nil {
		err = errs.Wrap(ferr, errs.WithContext("url", urlStr))
		return
	}
	link = &Webpage{
		URL:         urlStr,
		Canonical:   wi.Canonical,
		Title:       wi.Title,
		Description: wi.Description,
		ImageURL:    wi.ImageURL,
	}
	return
}

// SortPages function sorts Info list.
func SortPages(webpages []*Webpage) {
	if len(webpages) < 2 {
		return
	}
	sort.SliceStable(webpages, func(i, j int) bool {
		if webpages[i].Published == nil && webpages[j].Published == nil {
			return true
		}
		if webpages[i].Published == nil {
			return true
		}
		if webpages[j].Published == nil {
			return false
		}
		return webpages[i].Published.Before(*webpages[j].Published)
	})
}

// Encode putputs to io.Writer by JSON format.
func (i *Webpage) Encode(w io.Writer) error {
	if err := json.NewEncoder(w).Encode(i); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (wp *Webpage) ImageFile(ctx context.Context, dir string) (tname string, err error) {
	if wp == nil {
		err = errs.Wrap(ecode.ErrNullPointer)
		return
	}
	if len(wp.ImageURL) == 0 {
		err = errs.Wrap(ecode.ErrNoAPODImage)
		return
	}
	wi := &webinfo.Webinfo{ImageURL: wp.ImageURL, UserAgent: webinfo.DefaultUserAgent()}
	tname, err = wi.DownloadImage(ctx, dir, true)
	return
}

func (wp *Webpage) MakeMessage(prefixMsg string) string {
	if wp == nil {
		return ""
	}
	bld := strings.Builder{}

	//title
	if len(wp.Title) > 0 {
		fmt.Fprintln(&bld, prefixMsg, wp.Title)
	}
	// URL
	fmt.Fprintln(&bld, wp.URL)
	return bld.String()
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
