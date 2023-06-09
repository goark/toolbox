package webpage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goark/errs"
	"github.com/goark/fetch"
	"github.com/goark/toolbox/ecode"
	"github.com/mattn/go-encoding"
	"golang.org/x/net/html/charset"
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
func ReadPage(ctx context.Context, urlStr string) (*Webpage, error) {
	// fetch web page
	u, err := fetch.URL(urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	resp, err := fetch.New().GetWithContext(ctx, u)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", u.String()))
	}
	defer resp.Close()

	// detect character encoding
	br := bufio.NewReader(resp.Body())
	var r io.Reader = br
	if data, err2 := br.Peek(1024); err2 == nil { //next 1024 bytes without advancing the reader.
		enc, name, _ := charset.DetermineEncoding(data, resp.Header().Get("content-type"))
		if enc != nil {
			r = enc.NewDecoder().Reader(br)
		} else if len(name) > 0 {
			if enc := encoding.GetEncoding(name); enc != nil {
				r = enc.NewDecoder().Reader(br)
			}
		}
	}

	// analysis web content
	link := &Webpage{URL: urlStr}
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", u.String()))
	}
	doc.Find("head").Each(func(_ int, s *goquery.Selection) {
		s.Find("title").Each(func(_ int, s *goquery.Selection) {
			t := s.Text()
			if len(t) > 0 {
				link.Title = strings.TrimSpace(t)
			}
		})
		s.Find(`meta[property="og:title"]`).Each(func(_ int, s *goquery.Selection) {
			if v, ok := s.Attr("content"); ok && len(v) > 0 {
				link.Title = strings.TrimSpace(v)
			}
		})
		s.Find(`meta[name="description"]`).Each(func(_ int, s *goquery.Selection) {
			if v, ok := s.Attr("content"); ok && len(v) > 0 {
				link.Description = strings.TrimSpace(v)
			}
		})
		s.Find(`meta[property="og:description"]`).Each(func(_ int, s *goquery.Selection) {
			if v, ok := s.Attr("content"); ok && len(v) > 0 {
				link.Description = strings.TrimSpace(v)
			}
		})
		s.Find(`meta[property="og:image"]`).Each(func(_ int, s *goquery.Selection) {
			if v, ok := s.Attr("content"); ok && len(v) > 0 {
				link.ImageURL = strings.TrimSpace(v)
			}
		})
		s.Find("link[rel='canonical']").Each(func(_ int, s *goquery.Selection) {
			if v, ok := s.Attr("href"); ok && len(v) > 0 {
				link.Canonical = strings.TrimSpace(v)
			}
		})
	})
	return link, nil
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

func (wp *Webpage) ImageFile(ctx context.Context, dir string) (string, error) {
	if wp == nil {
		return "", errs.Wrap(ecode.ErrNullPointer)
	}
	if len(wp.ImageURL) == 0 {
		return "", errs.Wrap(ecode.ErrNoAPODImage)
	}

	// get Image data
	u, err := url.Parse(wp.ImageURL)
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("image_url", wp.ImageURL))
	}
	img, err := fetch.New().GetWithContext(ctx, u)
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("image_url", wp.ImageURL))
	}
	defer img.Close()

	// copy to temporary file
	file, err := os.CreateTemp(dir, "webpage.*.jpg")
	if err != nil {
		return "", errs.Wrap(err)
	}
	defer file.Close()

	tname := file.Name()
	_, err = io.Copy(file, img.Body())
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("image_url", wp.ImageURL), errs.WithContext("temp_file", tname))
	}
	return tname, nil
}

func (wp *Webpage) MakeMessage(prefixMsg string) string {
	if wp == nil {
		return ""
	}
	bld := strings.Builder{}

	//title
	if len(wp.Title) > 0 {
		bld.WriteString(fmt.Sprintln(prefixMsg, wp.Title))
	}
	// URL
	bld.WriteString(fmt.Sprintln(wp.URL))
	return bld.String()
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
