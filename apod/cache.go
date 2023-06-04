package apod

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/nasaapi/nasaapod"
	"github.com/goark/toolbox/values"
	"go.uber.org/zap"
)

const (
	cachesFile = "apod.cache.json"
)

type cacheData struct {
	Date   values.Date          `json:"date,omitempty"`
	Caches []*nasaapod.Response `json:"caches,omitempty"`
}

func (cfg *APOD) importCacheData() (*cacheData, error) {
	if cfg == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	path := filepath.Join(cfg.cacheDir, cachesFile)
	cfg.Logger().Debug("start importing cache file", zap.String("path", path))
	file, err := os.Open(path)
	if err != nil {
		cfg.Logger().Info("fail open cache file", zap.Object("error", zapobject.New(errs.Wrap(err, errs.WithContext("cache_dir", cfg.cacheDir), errs.WithContext("path", path)))))
		return &cacheData{}, nil
	}
	defer file.Close()

	var data cacheData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("cache_dir", cfg.cacheDir), errs.WithContext("path", path))
	}
	cfg.Logger().Debug("complete importing cache file", zap.String("path", path))
	data.sort()
	return &data, nil
}

func (cfg *APOD) exportCacheData(data *cacheData) error {
	if data == nil || len(data.Caches) == 0 {
		return nil
	}
	data.sort()
	if cfg == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	path := filepath.Join(cfg.cacheDir, cachesFile)
	cfg.Logger().Debug("start exporting cache file", zap.String("path", path))
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return errs.Wrap(err, errs.WithContext("cache_dir", cfg.cacheDir), errs.WithContext("path", path))
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(data); err != nil {
		return errs.Wrap(err, errs.WithContext("cache_dir", cfg.cacheDir), errs.WithContext("path", path))
	}
	cfg.Logger().Debug("complete exporting cache file", zap.String("path", path))
	return nil
}

func (c *cacheData) sort() {
	if c == nil || len(c.Caches) <= 1 {
		return
	}
	sort.Slice(c.Caches, func(i, j int) bool {
		return c.Caches[i].Date.Before(c.Caches[j].Date)
	})
}

func (c *cacheData) find(date values.Date) *nasaapod.Response {
	if c == nil || len(c.Caches) == 0 {
		return nil
	}
	i := sort.Search(len(c.Caches), func(i int) bool { return !c.Caches[i].Date.Before(date) })
	if i < len(c.Caches) && c.Caches[i].Date.Equal(date) {
		return c.Caches[i]
	}
	return nil
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
