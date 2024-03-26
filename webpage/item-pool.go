package webpage

import (
	"context"
	"strings"
	"sync"

	"github.com/goark/errs"
	"github.com/goark/toolbox/webpage/feed"
)

type itemPool struct {
	wg      sync.WaitGroup
	pool    *infoPool
	errList *errs.Errors
}

func newItemPool() *itemPool {
	pool := &itemPool{pool: newPool(), errList: &errs.Errors{}}
	pool.pool.start()
	return pool
}

func (ip *itemPool) putFeedItem(ctx context.Context, item *feed.Item) {
	if ip == nil {
		return
	}
	ip.wg.Add(1)
	go func() {
		defer ip.wg.Done()
		page, err := convWebpageFromFeedItem(ctx, item)
		if err != nil {
			ip.errList.Add(err)
			return
		}
		ip.putPage(page)
	}()
}

func (ip *itemPool) putPage(page *Webpage) {
	if ip == nil {
		return
	}
	ip.pool.put(page)
}

func (ip *itemPool) done() {
	if ip == nil {
		return
	}
	ip.wg.Wait()
	ip.pool.stop()
}

func (ip *itemPool) length() int {
	if ip == nil {
		return 0
	}
	return ip.pool.length()
}

func (ip *itemPool) getPages() []*Webpage {
	if ip == nil {
		return []*Webpage{}
	}
	return ip.pool.getList()
}

const (
	githubDomainInURL = "//github.com/"
)

func convWebpageFromFeedItem(ctx context.Context, item *feed.Item) (*Webpage, error) {
	page := &Webpage{
		URL:         item.Link,
		Title:       item.Title,
		Description: item.Description,
		Published:   item.Published,
	}
	if len(item.Images) > 0 {
		page.ImageURL = item.Images[0].URL
	}
	if len(page.ImageURL) == 0 || strings.Contains(item.Link, githubDomainInURL) {
		i, err := ReadPage(ctx, page.URL)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("url", page.URL))
		}
		if strings.Contains(item.Link, githubDomainInURL) {
			page.Title = i.Title
		}
		if len(page.ImageURL) == 0 {
			page.ImageURL = i.ImageURL
		}
	}
	return page, nil
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
