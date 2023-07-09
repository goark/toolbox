package apod

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/toolbox/db/model"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/nasaapi/nasaapod"
	"github.com/goark/toolbox/values"
	"go.uber.org/zap"
)

func (cfg *APOD) find(ctx context.Context, date values.Date) (*nasaapod.Response, error) {
	if cfg == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	if data, ok := cfg.cache[date.String()]; ok {
		return data, nil
	}
	data, err := cfg.repos.FindAPODDataByDate(ctx, date.String())
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("date", date.String()))
	} else if data == nil {
		cfg.Logger().Debug("no data in database", zap.String("date", date.String()))
		return nil, nil
	}
	apoddata := importFromModel(data)
	cfg.cache[date.String()] = apoddata
	return apoddata, nil
}

func (cfg *APOD) saveDB(ctx context.Context) error {
	if cfg == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	if len(cfg.saveData) == 0 {
		return nil
	}
	data := make([]model.ApodData, 0, len(cfg.saveData))
	for _, d := range cfg.saveData {
		data = append(data, exportToModel(d))
	}
	cfg.Logger().Debug("start saving data to database", zap.Any("data", data))
	if err := cfg.repos.InsertAPODData(ctx, data); err != nil {
		return errs.Wrap(err)
	}
	cfg.Logger().Debug("complete saving data to database", zap.Any("data", data))
	return nil
}

func (cfg *APOD) put(data *nasaapod.Response) {
	cfg.saveData = append(cfg.saveData, data)
}

func importFromModel(data *model.ApodData) *nasaapod.Response {
	if data == nil {
		return nil
	}
	date, _ := values.DateFrom(data.Date, false)
	return &nasaapod.Response{
		Copyright:      data.Copyright,
		Date:           date,
		Explanation:    data.Explanation,
		HdUrl:          data.HdUrl,
		MediaType:      data.MediaType,
		ServiceVersion: data.ServiceVersion,
		Title:          data.Title,
		Url:            data.Url,
		ThumbnailUrl:   data.ThumbnailUrl,
	}
}

func exportToModel(data *nasaapod.Response) model.ApodData {
	if data == nil {
		return model.ApodData{}
	}
	return model.ApodData{
		Copyright:      data.Copyright,
		Date:           data.Date.String(),
		Explanation:    data.Explanation,
		HdUrl:          data.HdUrl,
		MediaType:      data.MediaType,
		ServiceVersion: data.ServiceVersion,
		Title:          data.Title,
		Url:            data.Url,
		ThumbnailUrl:   data.ThumbnailUrl,
	}
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
