package apod

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/nasaapi/nasaapod"
	"github.com/goark/toolbox/values"
	"go.uber.org/zap"
)

func (cfg *APOD) Lookup(ctx context.Context, date values.Date) (*nasaapod.Response, error) {
	if cfg == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	if date.IsZero() {
		date = values.Today()
	}
	dt, err := cfg.importCacheData()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if !dt.Date.IsZero() && !dt.Date.Before(date) {
		if res := dt.find(date); res != nil {
			cfg.Logger().Debug("find data in cache", zap.Any("data", res))
			return res, nil
		}
	}

	cfg.Logger().Debug("start reading APOD data", zap.String("date", date.String()))
	res, err := nasaapod.New(
		nasaapod.WithAPIKey(cfg.APIKey),
		nasaapod.WithDate(date),
		nasaapod.WithThumbs(true),
	).Get(ctx)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("date", date.String()))
	}
	cfg.Logger().Debug("complete reading APOD data", zap.String("date", date.String()), zap.Any("response", res))
	if len(res) == 0 {
		return nil, errs.Wrap(ecode.ErrNoContent, errs.WithContext("date", date.String()))
	}
	dt.Caches = append(dt.Caches, res...)
	dt.Date = values.Today()
	cfg.Logger().Debug("cache data", zap.Any("data", dt))
	if err := cfg.exportCacheData(dt); err != nil {
		return nil, errs.Wrap(ecode.ErrNoContent, errs.WithContext("date", date.String()))
	}
	return res[0], nil
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
