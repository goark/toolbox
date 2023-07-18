package webpage

import (
	"context"
	"net/url"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/webpage/feed"
	"go.uber.org/zap"
)

// Feed fetches feed URL and gets webpage informations.
func (cfg *Config) Feed(ctx context.Context, urlStr string) error {
	if cfg == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	resp, err := inmortFeed(ctx, urlStr)
	if err != nil {
		return errs.Wrap(err, errs.WithContext("feed_url", urlStr))
	}
	cfg.getNewDataList(ctx, resp)
	return nil
}

// Feed fetches feed URL and gets webpage informations.
func (cfg *Config) FeedFlickr(ctx context.Context, flickrId string) error {
	if cfg == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	resp, err := inmortFeedFlickr(ctx, flickrId)
	if err != nil {
		return errs.Wrap(err, errs.WithContext("flickr_id", flickrId))
	}
	cfg.getNewDataList(ctx, resp)
	return nil
}

func (cfg *Config) getNewDataList(ctx context.Context, items []*feed.Item) {
	if cfg.itemPool == nil {
		cfg.CreatePool()
	}
	urls := map[string]bool{}
	for _, item := range items {
		if urls[item.Link] {
			continue
		}
		urls[item.Link] = true
		if page, err := cfg.find(ctx, item.Link); err != nil || page == nil {
			cfg.itemPool.putFeedItem(ctx, item)
			cfg.Logger().Debug("new item", zap.Any("item", item))
		}
	}
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
