package apod

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

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
func New(ctx context.Context, path, cacheDir string, logger *log.ZapEventLogger) (cfg *APOD, err error) {
	// open database
	repos, ferr := db.Open(ctx, cacheDir, logger)
	if ferr != nil {
		err = errs.Wrap(ferr, errs.WithContext("cache_dir", cacheDir))
		return
	}

	// read configuration file
	if len(path) == 0 {
		cfg = fallthroughCfg(repos, logger)
		return
	}
	file, ferr := os.Open(filepath.Clean(path))
	if ferr != nil {
		cfg = fallthroughCfg(repos, logger)
		return
	}
	defer func() {
		err = errs.Join(err, file.Close())
	}()
	var vcfg APOD
	if jerr := json.NewDecoder(file).Decode(&vcfg); jerr != nil {
		err = errs.Wrap(jerr, errs.WithContext("path", path))
		return
	}
	vcfg.logger = logger
	vcfg.cacheDir = cacheDir
	vcfg.repos = repos
	vcfg.cache = map[string]*nasaapod.Response{}
	vcfg.saveData = []*nasaapod.Response{}
	cfg = &vcfg
	return
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
func (cfg *APOD) Export(path string) (err error) {
	if cfg == nil {
		err = errs.Wrap(ecode.ErrNullPointer)
		return
	}
	file, ferr := os.OpenFile(filepath.Clean(path), os.O_RDWR|os.O_CREATE, 0600)
	if ferr != nil {
		err = errs.Wrap(ferr, errs.WithContext("path", path))
		return
	}
	defer func() {
		err = errs.Join(err, file.Close())
	}()

	if jerr := json.NewEncoder(file).Encode(cfg); jerr != nil {
		err = errs.Wrap(jerr, errs.WithContext("path", path))
		return
	}
	return
}

/* Copyright 2023-2026 Spiegel
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
