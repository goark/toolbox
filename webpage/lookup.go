package webpage

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"go.uber.org/zap"
)

func (cfg *Config) Lookup(ctx context.Context, urlStr string) (*Webpage, error) {
	page, exist, err := cfg.GetWebpage(ctx, urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	if !exist {
		cfg.Logger().Debug("put webpage data to pool", zap.Any("webpage", page))
		cfg.itemPool.putPage(page)
	}
	return page, nil
}

func (cfg *Config) GetWebpage(ctx context.Context, urlStr string) (*Webpage, bool, error) {
	if cfg == nil {
		return nil, false, errs.Wrap(ecode.ErrNullPointer)
	}
	// find data from cache or database
	if page, err := cfg.find(ctx, urlStr); err != nil {
		return nil, false, errs.Wrap(err, errs.WithContext("url", urlStr))
	} else if page != nil {
		cfg.Logger().Debug("get webpage data from database", zap.Any("webpage", page))
		return page, true, nil
	}
	// fetch webpage.
	page, err := ReadPage(ctx, urlStr)
	if err != nil {
		return nil, false, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	cfg.Logger().Debug("fetch webpage data", zap.Any("webpage", page))
	return page, false, nil
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
