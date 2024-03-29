package db

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/toolbox/db/model"
	"github.com/goark/toolbox/ecode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FindWebpageByDate method finds Webpage data from database condition by url.
func (repos *Repository) FindWebpageByURL(ctx context.Context, url string) (*model.Webpage, error) {
	if repos == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	var data model.Webpage
	tx := repos.Db().WithContext(ctx).Where(&model.Webpage{URL: url}).Order("created_at asc").First(&data)
	if tx.Error != nil {
		err := errs.Wrap(tx.Error, errs.WithContext("url", url))
		if errs.Is(tx.Error, gorm.ErrRecordNotFound) {
			repos.Logger().Debug("no record", zap.Object("error", zapobject.New(err)))
			return nil, nil
		}
		return nil, err
	}
	repos.Logger().Debug("find data", zap.Any("data", data))
	return &data, nil
}

// InsertAPODData method inserts APOD data to database.
func (repos *Repository) InsertWebpage(ctx context.Context, datalist []model.Webpage) error {
	if repos == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	for _, data := range datalist {
		data := data
		if err := repos.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if t := tx.Create(&data); t.Error != nil {
				return errs.Wrap(t.Error)
			}
			return nil
		}); err != nil {
			return err
		}
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
