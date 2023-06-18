package feed

import (
	"context"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/goark/errs"
	ftch "github.com/goark/fetch"
	"github.com/mmcdole/gofeed/atom"
)

// FeedFlickr fetches feed data from Flickr ID.
func FeedFlickr(ctx context.Context, flickrId string) (*Metadata, error) {
	u, err := makeFlickrFeedURL(flickrId)
	if err != nil {
		return nil, errs.Wrap(ErrInvalidFlickrId, errs.WithContext("flickr_id", flickrId))
	}
	resp, err := ftch.New().GetWithContext(ctx, u)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", u.String()))
	}
	defer resp.Close()

	return decodeFlickrFeed(resp.Body())
}

func makeFlickrFeedURL(flickrId string) (*url.URL, error) {
	u, err := url.Parse("https://www.flickr.com/services/feeds/photos_public.gne?format=atom")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if len(flickrId) == 0 {
		return nil, errs.Wrap(ErrInvalidFlickrId, errs.WithContext("flickr_id", flickrId))
	}
	q := u.Query()
	q.Add("id", flickrId)
	u.RawQuery = q.Encode()
	return u, nil
}

func decodeFlickrFeed(r io.Reader) (*Metadata, error) {
	f, err := (&atom.Parser{}).Parse(r)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	data := &Metadata{
		Title:       f.Title,
		Description: f.Subtitle,
	}
	for _, lnk := range f.Links {
		switch {
		case strings.EqualFold(lnk.Rel, "self"):
			data.FeedLink = lnk.Href
		case strings.EqualFold(lnk.Rel, "alternate"):
			data.Link = lnk.Href
		}
	}
	items := []*Item{}
	for _, i := range f.Entries {
		item := &Item{
			Title:     i.Title,
			Published: i.PublishedParsed,
			Updated:   i.UpdatedParsed,
		}
		image := &Image{}
		for _, lnk := range i.Links {
			switch {
			case strings.EqualFold(lnk.Rel, "alternate"):
				item.Link = lnk.Href
			case strings.EqualFold(lnk.Rel, "enclosure"):
				image.MimeType = lnk.Type
				image.URL = lnk.Href
			}
		}
		if i.Extensions != nil {
			if flickr, ok := i.Extensions["flickr"]; ok {
				if taken, ok := flickr["date_taken"]; ok && len(taken) > 0 {
					if parsedTaken, err := time.Parse(time.RFC3339, taken[0].Value); err == nil {
						parsedTaken = parsedTaken.In(time.UTC)
						image.Taken = &parsedTaken
					}
				}
			}
		}
		item.Images = []*Image{image}
		authors := []*Author{}
		for _, a := range i.Authors {
			authors = append(authors, &Author{Name: a.Name, URL: a.URI})
		}
		item.Authors = authors
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
