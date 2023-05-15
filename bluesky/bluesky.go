package bluesky

import (
	"encoding/json"
	"os"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/rs/zerolog"
)

const (
	DefaltHostName = "bsky.social"
)

// Bluesky is configuration for Bluesky
type Bluesky struct {
	Host     string `json:"host"`
	Handle   string `json:"handle"`
	Password string `json:"password"`
	baseDir  string
	logger   *zerolog.Logger
	client   *xrpc.Client
}

// New creates new Bluesky instance.
func New(path, dir string, logger *zerolog.Logger) (*Bluesky, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path), errs.WithContext("die", dir))
	}
	defer file.Close()

	var cfg Bluesky
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	if len(cfg.Host) == 0 {
		cfg.Host = "https://" + DefaltHostName
	}
	if len(cfg.Handle) == 0 {
		return nil, errs.Wrap(ecode.ErrNoBlueskyHandle, errs.WithContext("path", path), errs.WithContext("die", dir))
	}
	cfg.baseDir = dir
	if logger == nil {
		lggr := zerolog.Nop()
		logger = &lggr
	}
	cfg.logger = logger
	return &cfg, nil
}

// BaseDir method returns base directory.
func (cfg *Bluesky) BaseDir() string {
	if cfg == nil {
		return ""
	}
	return cfg.baseDir
}

// Logger method returns zerolog.Logger instance.
func (cfg *Bluesky) Logger() *zerolog.Logger {
	if cfg == nil || cfg.logger == nil {
		logger := zerolog.Nop()
		return &logger
	}
	return cfg.logger
}

// Export methods exports configuration to config file.
func (cfg *Bluesky) Export(path string) error {
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
