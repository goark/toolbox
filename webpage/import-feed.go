package webpage

import (
	"context"
	"net/url"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/webpage/feed"
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
	return wp.getNewDataList(resp), nil
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
	return wp.getNewDataList(resp), nil
}

func (wp *Webpage) getNewDataList(infos []*Info) []*Info {
	list := []*Info{}
	for _, info := range infos {
		if i := wp.cacheData.Get(info.URL); i == nil {
			list = append(list, info)
			wp.cacheData.Put(info)
		}
	}
	return list
}

func inmortFeed(ctx context.Context, urlStr string) ([]*Info, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	data, err := feed.Feed(ctx, u)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	list := []*Info{}
	for _, item := range data.Items {
		info, err := importWebpage(ctx, item)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		list = append(list, info)
	}
	return list, nil
}

func inmortFeedFlickr(ctx context.Context, flickrId string) ([]*Info, error) {
	data, err := feed.FeedFlickr(ctx, flickrId)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("flickr_id", flickrId))
	}
	list := []*Info{}
	for _, item := range data.Items {
		info, err := importWebpage(ctx, item)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		list = append(list, info)
	}
	return list, nil
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
	if len(info.ImageURL) == 0 {
		i, err := ReadPage(ctx, info.URL)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("url", info.URL))
		}
		info.ImageURL = i.ImageURL
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
