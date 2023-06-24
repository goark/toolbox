package webpage

import (
	"context"
	"net/url"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/webpage/feed"
	"go.uber.org/zap"
)

const (
	githubDomainInURL = "//github.com/"
)

// Feed fetches feed URL and gets webpage informations.
func (wp *Webpage) Feed(ctx context.Context, urlStr string) ([]*Info, error) {
	if wp == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	resp, err := inmortFeed(ctx, urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("feed_url", urlStr))
	}
	return wp.getNewDataList(ctx, resp)
}

// Feed fetches feed URL and gets webpage informations.
func (wp *Webpage) FeedFlickr(ctx context.Context, flickrId string) ([]*Info, error) {
	if wp == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	resp, err := inmortFeedFlickr(ctx, flickrId)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("flickr_id", flickrId))
	}
	return wp.getNewDataList(ctx, resp)
}

func (wp *Webpage) getNewDataList(ctx context.Context, items []*feed.Item) ([]*Info, error) {
	list := []*Info{}
	for _, item := range items {
		if i := wp.cacheData.Get(item.Link); i == nil {
			info, err := importWebpage(ctx, item)
			if err != nil {
				return nil, errs.Wrap(err, errs.WithContext("flickr_id", errs.WithContext("url", item.Link)))
			}
			list = append(list, info)
			wp.cacheData.Put(info)
			wp.Logger().Debug("put web page to cache", zap.Any("info", info))
		}
	}
	return list, nil
}

func inmortFeed(ctx context.Context, urlStr string) ([]*feed.Item, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	data, err := feed.Feed(ctx, u)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	if data == nil || data.Items == nil {
		return []*feed.Item{}, nil
	}
	return data.Items, nil
}

func inmortFeedFlickr(ctx context.Context, flickrId string) ([]*feed.Item, error) {
	data, err := feed.FeedFlickr(ctx, flickrId)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("flickr_id", flickrId))
	}
	if data == nil || data.Items == nil {
		return []*feed.Item{}, nil
	}
	return data.Items, nil
}

func importWebpage(ctx context.Context, item *feed.Item) (*Info, error) {
	info := &Info{
		URL:         item.Link,
		Title:       item.Title,
		Description: item.Description,
		Published:   item.Published,
	}
	if len(item.Images) > 0 {
		info.ImageURL = item.Images[0].URL
	}
	if len(info.ImageURL) == 0 || strings.Contains(item.Link, githubDomainInURL) {
		i, err := ReadPage(ctx, info.URL)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("url", info.URL))
		}
		if strings.Contains(item.Link, githubDomainInURL) {
			info.Title = i.Title
		}
		if len(info.ImageURL) == 0 {
			info.ImageURL = i.ImageURL
		}
	}
	return info, nil
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
