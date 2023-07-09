package conn

import (
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/goark/errs"
	"github.com/ipfs/go-log/v2"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	connDB  *gorm.DB
	errorDB error
	onceDB  sync.Once
)

func GetDB(path string, zlogger *log.ZapEventLogger) (*gorm.DB, error) {
	onceDB.Do(func() {
		conn, err := openDB(path, zlogger)
		if err != nil {
			errorDB = errs.Wrap(err, errs.WithContext("dbfile", path))
		} else {
			connDB = conn
		}
	})
	return connDB, errorDB
}

func openDB(path string, zlogger *log.ZapEventLogger) (*gorm.DB, error) {
	// logger
	gcfg := &gorm.Config{}
	if lggr, lvl := getLogger(zlogger); lvl != gormlogger.Silent {
		gcfg.Logger = lggr
	}
	// open SQLite database
	db, err := gorm.Open(sqlite.Open(path), gcfg)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("dbfile", path))
	}
	return db, nil
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
