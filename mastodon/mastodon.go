package mastodon

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/goark/errs"
	"github.com/goark/toolbox/consts"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/logger"
	"github.com/ipfs/go-log/v2"
	mstdn "github.com/mattn/go-mastodon"
	"go.uber.org/zap"
)

const (
	scopes = "read write follow"
)

// Mastodon is configuration for Mastodon
type Mastodon struct {
	Server       string `json:"server"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	client       *mstdn.Client
	logger       *log.ZapEventLogger
}

func New(path string, logger *log.ZapEventLogger) (*Mastodon, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	defer file.Close()

	var cfg Mastodon
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", path))
	}
	cfg.client = mstdn.NewClient(&mstdn.Config{
		Server:       cfg.Server,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		AccessToken:  cfg.AccessToken,
	})
	cfg.logger = logger
	return &cfg, nil
}

// AppName method returns application name.
func (cfg *Mastodon) AppName() string {
	return consts.AppName
}

// Scopes method returns scopes of application.
func (cfg *Mastodon) Scopes() string {
	return scopes
}

// Registory method returns registory URL of application.
func (cfg *Mastodon) Registory() string {
	return consts.RepositoryURL
}

func (cfg *Mastodon) Servername() string {
	if cfg == nil {
		return ""
	}
	if u, err := url.Parse(cfg.Server); err == nil && len(u.Hostname()) > 0 {
		return u.Hostname()
	}
	return cfg.Server
}

// Logger method returns zap.Logger instance.
func (cfg *Mastodon) Logger() *zap.Logger {
	if cfg == nil || cfg.logger == nil {
		return logger.Nop().Desugar()
	}
	return cfg.logger.Desugar()
}

// Export methods exports configuration to config file.
func (cfg *Mastodon) Export(path string) error {
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
