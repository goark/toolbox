package db

import (
	"context"
	"os"
	"path/filepath"

	"github.com/goark/errs"
	"github.com/goark/toolbox/db/conn"
	"github.com/goark/toolbox/db/model"
	"github.com/goark/toolbox/logger"
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	sqliteFile = "db.sqlite"
)

func existFile(dir string) (string, bool) {
	path := filepath.Join(dir, sqliteFile)
	if _, err := os.Stat(path); err != nil {
		return path, false
	}
	return path, true
}

type Repository struct {
	db     *gorm.DB
	logger *log.ZapEventLogger
}

func Open(ctx context.Context, dir string, zlogger *log.ZapEventLogger) (*Repository, error) {
	// path of SQLite file
	path, existFlag := existFile(dir)
	zlogger.Desugar().Debug("database file", zap.String("path", path), zap.Bool("file exist", existFlag))
	// open SQLite database
	db, err := conn.Open(path, zlogger)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("dbfile", path))
	}
	zlogger.Desugar().Debug("complete opening database file", zap.String("path", path))
	// migration
	if !existFlag {
		zlogger.Desugar().Debug("start migration", zap.String("path", path), zap.Bool("file exist", existFlag))
		if err := model.Migration(ctx, db); err != nil {
			return nil, errs.Wrap(err, errs.WithContext("dbfile", path))
		}
		zlogger.Desugar().Debug("complete migration", zap.String("path", path), zap.Bool("file exist", existFlag))
	}

	return &Repository{db: db, logger: zlogger}, nil
}

// Db method returns gorm.DB instance.
func (cfg *Repository) Db() *gorm.DB {
	if cfg == nil {
		return nil
	}
	return cfg.db
}

// Logger method returns zap.Logger instance.
func (cfg *Repository) Logger() *zap.Logger {
	if cfg == nil || cfg.logger == nil {
		return logger.Nop().Desugar()
	}
	return cfg.logger.Desugar()
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
