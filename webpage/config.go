package webpage

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/toolbox/db"
	"github.com/goark/toolbox/logger"
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
)

// Config is configuration for webpage
type Config struct {
	cacheDir  string
	cacheData *Cache
	itemPool  *itemPool
	logger    *log.ZapEventLogger
	repos     *db.Repository
}

// New functions creates new Config instance.
func New(ctx context.Context, cacheDir string, logger *log.ZapEventLogger) (*Config, error) {
	// open database
	repos, err := db.Open(ctx, cacheDir, logger)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("cache_dir", cacheDir))
	}
	// make configuration
	cfg := &Config{
		cacheDir:  cacheDir,
		cacheData: NewCache(cacheDir),
		logger:    logger,
		repos:     repos,
	}
	cfg.CreatePool()
	return cfg, nil
}

// Logger method returns zap.Logger instance.
func (cfg *Config) Logger() *zap.Logger {
	if cfg == nil || cfg.logger == nil {
		return logger.Nop().Desugar()
	}
	return cfg.logger.Desugar()
}

func (cfg *Config) Save(ctx context.Context, list []*Webpage) error {
	if len(list) == 0 {
		cfg.Logger().Debug("no save data in pool")
		return nil
	}
	return cfg.saveDB(ctx, list)
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
