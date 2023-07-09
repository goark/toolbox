package webpage

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/toolbox/db/model"
	"github.com/goark/toolbox/ecode"
	"go.uber.org/zap"
)

func (cfg *Config) find(ctx context.Context, url string) (*Webpage, error) {
	if cfg == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	// get webpage info from cache
	if page := cfg.cacheData.Get(url); page != nil {
		return page, nil
	}
	data, err := cfg.repos.FindWebpageByURL(ctx, url)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", url))
	} else if data == nil {
		cfg.Logger().Debug("no data in database", zap.String("url", url))
		return nil, nil
	}
	page := importWebpageFromModel(*data)
	cfg.cacheData.Put(page)
	return page, nil
}

func (cfg *Config) saveDB(ctx context.Context, list []*Webpage) error {
	// make save data
	data := make([]model.Webpage, 0, len(list))
	for _, d := range list {
		data = append(data, exportWebpageToModel(d))
	}
	cfg.Logger().Debug("start saving data to database", zap.Any("data", data))
	if err := cfg.repos.InsertWebpage(ctx, data); err != nil {
		return errs.Wrap(err)
	}
	cfg.Logger().Debug("complete saving data to database")
	return nil
}

func importWebpageFromModel(data model.Webpage) *Webpage {
	return &Webpage{
		URL:         data.URL,
		Canonical:   data.Canonical,
		Title:       data.Title,
		Description: data.Description,
		ImageURL:    data.ImageURL,
		Published:   data.GetPublished(),
	}
}

func exportWebpageToModel(page *Webpage) model.Webpage {
	if page == nil {
		return model.Webpage{}
	}
	data := model.Webpage{
		URL:         page.URL,
		Canonical:   page.Canonical,
		Title:       page.Title,
		Description: page.Description,
		ImageURL:    page.ImageURL,
	}
	data.SetPublished(page.Published)
	return data
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
