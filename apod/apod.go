package apod

import (
	"context"
	"encoding/json"
	"os"

	"github.com/goark/errs"
	"github.com/goark/toolbox/db"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/logger"
	"github.com/goark/toolbox/nasaapi"
	"github.com/goark/toolbox/nasaapi/nasaapod"
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
)

// APOD is configuration for NASA API and APOD
type APOD struct {
	APIKey   string `json:"api_key"`
	cacheDir string
	logger   *log.ZapEventLogger
	repos    *db.Repository
	cache    map[string]*nasaapod.Response
	saveData []*nasaapod.Response
}

// New functions creates new APOD instance from file.
func New(ctx context.Context, path, cacheDir string, logger *log.ZapEventLogger) (*APOD, error) {
	// open database
	repos, err := db.Open(ctx, cacheDir, logger)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("cache_dir", cacheDir))
	}

	// read configuration file
	if len(path) == 0 {
		return fallthroughCfg(repos, logger), nil
	}
	file, err := os.Open(path)
	if err != nil {
		return fallthroughCfg(repos, logger), nil
	}
	defer file.Close()
	var cfg APOD
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	cfg.logger = logger
	cfg.cacheDir = cacheDir
	cfg.repos = repos
	cfg.cache = map[string]*nasaapod.Response{}
	cfg.saveData = []*nasaapod.Response{}

	return &cfg, nil
}

func fallthroughCfg(repos *db.Repository, logger *log.ZapEventLogger) *APOD {
	return &APOD{
		APIKey:   nasaapi.DefaultAPIKey,
		logger:   logger,
		repos:    repos,
		cache:    map[string]*nasaapod.Response{},
		saveData: []*nasaapod.Response{},
	}
}

// Logger method returns zap.Logger instance.
func (cfg *APOD) Logger() *zap.Logger {
	if cfg == nil || cfg.logger == nil {
		return logger.Nop().Desugar()
	}
	return cfg.logger.Desugar()
}

// Export methods exports configuration to config file.
func (cfg *APOD) Export(path string) error {
	if cfg == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return errs.Wrap(err, errs.WithContext("path", path))
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(cfg); err != nil {
		return errs.Wrap(err, errs.WithContext("path", path))
	}
	return nil
}

/* Copyright 2023-2024 Spiegel
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
