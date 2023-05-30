package apod

import (
	"encoding/json"
	"os"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/logger"
	"github.com/goark/toolbox/nasaapi"
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
)

// APOD is configuration for NASA API and APOD
type APOD struct {
	APIKey   string `json:"api_key"`
	cacheDir string
	logger   *log.ZapEventLogger
}

// New functions creates new APOD instance from file.
func New(path, cacheDir string, logger *log.ZapEventLogger) (*APOD, error) {
	if len(path) == 0 {
		return fallthroughCfg(cacheDir, logger), nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fallthroughCfg(cacheDir, logger), nil
	}
	defer file.Close()

	var cfg APOD
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	cfg.logger = logger
	cfg.cacheDir = cacheDir

	return &cfg, nil
}

func fallthroughCfg(cacheDir string, logger *log.ZapEventLogger) *APOD {
	return &APOD{
		APIKey: nasaapi.DefaultAPIKey,
		logger: logger,
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
