package feed

import (
	"context"
	"io"
	"net/url"

	"github.com/goark/errs"
	ftch "github.com/goark/fetch"
	"github.com/mmcdole/gofeed"
)

// Feed fetches feed data from URL.
func Feed(ctx context.Context, u *url.URL) (*Metadata, error) {
	resp, err := ftch.New().GetWithContext(ctx, u)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", u.String()))
	}
	defer resp.Close()

	return decodeFeed(resp.Body())
}

func decodeFeed(r io.Reader) (*Metadata, error) {
	f, err := gofeed.NewParser().Parse(r)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	data := &Metadata{
		FeedLink:    f.FeedLink,
		Title:       f.Title,
		Description: f.Description,
		Link:        f.Link,
	}
	authors := []*Author{}
	for _, a := range f.Authors {
		authors = append(authors, &Author{Name: a.Name})
	}
	data.Authors = authors
	items := []*Item{}
	for _, i := range f.Items {
		item := &Item{
			Title:       i.Title,
			Description: i.Description,
			Link:        i.Link,
			Published:   i.PublishedParsed,
			Updated:     i.UpdatedParsed,
		}
		authors := []*Author{}
		for _, a := range i.Authors {
			authors = append(authors, &Author{Name: a.Name})
		}
		item.Authors = authors
		if i.Image != nil {
			item.Images = []*Image{{Title: i.Image.Title, URL: i.Image.URL}}
		}
		items = append(items, item)
	}
	data.Items = items
	return data, nil
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
