package bluesky

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/logger"
	"github.com/goark/toolbox/webpage"
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
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
	wcfg     *webpage.Config
	logger   *log.ZapEventLogger
	client   *xrpc.Client
}

// New creates new Bluesky instance.
func New(path, dir string, wcfg *webpage.Config, logger *log.ZapEventLogger) (cfg *Bluesky, err error) {
	file, ferr := os.Open(filepath.Clean(path))
	if ferr != nil {
		err = errs.Wrap(ferr, errs.WithContext("path", path), errs.WithContext("die", dir))
		return
	}
	defer func() {
		err = errs.Join(err, file.Close())
	}()

	var vcfg Bluesky
	if jerr := json.NewDecoder(file).Decode(&vcfg); jerr != nil {
		err = errs.Wrap(err, errs.WithContext("path", path))
		return
	}
	if len(vcfg.Host) == 0 {
		cfg.Host = "https://" + DefaltHostName
	}
	if len(vcfg.Handle) == 0 {
		err = errs.Wrap(ecode.ErrNoBlueskyHandle, errs.WithContext("path", path), errs.WithContext("die", dir))
		return
	}
	vcfg.baseDir = dir
	vcfg.wcfg = wcfg
	vcfg.logger = logger
	cfg = &vcfg
	return
}

// BaseDir method returns base directory.
func (cfg *Bluesky) BaseDir() string {
	if cfg == nil {
		return ""
	}
	return cfg.baseDir
}

// Logger method returns zap.Logger instance.
func (cfg *Bluesky) Logger() *zap.Logger {
	if cfg == nil || cfg.logger == nil {
		return logger.Nop().Desugar()
	}
	return cfg.logger.Desugar()
}

// Export methods exports configuration to config file.
func (cfg *Bluesky) Export(path string) (err error) {
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
