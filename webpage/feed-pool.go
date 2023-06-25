package webpage

import (
	"context"
	"strings"
	"sync"

	"github.com/goark/errs"
	"github.com/goark/toolbox/webpage/feed"
)

type infoPool struct {
	mu   sync.RWMutex
	ch   chan *Info
	done chan struct{}
	pool []*Info
}

func newPool() *infoPool {
	return &infoPool{ch: make(chan *Info), done: make(chan struct{}), pool: []*Info{}}
}

func (p *infoPool) put(i *Info) {
	if p == nil {
		return
	}
	p.ch <- i
}

func (p *infoPool) start() {
	if p == nil {
		return
	}
	go func() {
		for {
			i, ok := <-p.ch
			if !ok {
				break
			}
			func() {
				p.mu.Lock()
				defer p.mu.Unlock()
				p.pool = append(p.pool, i)
			}()
		}
		p.done <- struct{}{}
	}()
}

func (p *infoPool) getList() []*Info {
	if p == nil {
		return []*Info{}
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.pool
}

func (p *infoPool) stop() {
	if p == nil {
		return
	}
	close(p.ch)
	<-p.done
}

type itemPool struct {
	ch      chan *feed.Item
	wg      sync.WaitGroup
	pool    *infoPool
	errList *errs.Errors
}

func newItemPool() *itemPool {
	pool := &itemPool{ch: make(chan *feed.Item), pool: newPool(), errList: &errs.Errors{}}
	pool.pool.start()
	return pool
}

func (ip *itemPool) put(ctx context.Context, item *feed.Item) {
	if ip == nil {
		return
	}
	ip.wg.Add(1)
	go func() {
		defer ip.wg.Done()
		info, err := convWebpageInfo(ctx, item)
		if err != nil {
			ip.errList.Add(err)
			return
		}
		ip.pool.put(info)
	}()
}

func (ip *itemPool) done() {
	if ip == nil {
		return
	}
	ip.wg.Wait()
	ip.pool.stop()
}

func (ip *itemPool) getInfo() []*Info {
	if ip == nil {
		return []*Info{}
	}
	return ip.pool.getList()
}

const (
	githubDomainInURL = "//github.com/"
)

func convWebpageInfo(ctx context.Context, item *feed.Item) (*Info, error) {
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
