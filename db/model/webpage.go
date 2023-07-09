package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Webpage struct {
	gorm.Model
	URL         string `gorm:"unique"`
	Canonical   string
	Title       string
	Description string
	ImageURL    string
	Published   sql.NullTime
}

// GetPublished returns pointer of time.Time for Webpage.Published.
func (wp Webpage) GetPublished() *time.Time {
	if !wp.Published.Valid {
		return nil
	}
	tm := wp.Published.Time
	return &tm
}

func (wp *Webpage) SetPublished(tm *time.Time) {
	if tm != nil {
		wp.Published = sql.NullTime{Time: *tm, Valid: true}
	} else {
		wp.Published = sql.NullTime{Time: time.Time{}, Valid: false}
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
