package webpage

import (
	"bufio"
	"context"
	"os"
	"strings"
	"sync"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
)

// FeedList is list of feed URLs.
type FeedList []string

// NewFeedList function returns new instance of FeedList.
func NewFeedList(path string) (FeedList, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	defer file.Close()

	list := FeedList{}
	s := bufio.NewScanner(file)
	for s.Scan() {
		txt := strings.TrimSpace(s.Text())
		if len(txt) > 0 {
			list = append(list, txt)
		}
	}
	if err := s.Err(); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	return list, nil
}

// Parse method parses feeds.
func (fl FeedList) Parse(ctx context.Context, wp *Webpage) error {
	if wp == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	var wg sync.WaitGroup
	errList := newErrorList()
	for _, urlStr := range fl {
		urlStr := urlStr
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := wp.Feed(ctx, urlStr); err != nil {
				errList.Add(errs.Wrap(err, errs.WithContext("feed_url", urlStr)))
			}
		}()
	}
	wg.Wait()
	return errList.GetError()
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
