package mastodon

import (
	"context"
	"net/url"

	"github.com/goark/errs"
	mstdn "github.com/mattn/go-mastodon"
	"github.com/rs/zerolog"
)

// Register functions registers application to mastodon server.
func Register(ctx context.Context, server, userId, password string, logger *zerolog.Logger) (*Mastodon, error) {
	if u, err := url.Parse(server); err == nil && len(u.Hostname()) > 0 {
		server = u.Hostname()
	}
	cfg := &Mastodon{
		Server: "https://" + server,
		logger: logger,
	}
	if err := cfg.register(ctx); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("server", cfg.Server))
	}
	if err := cfg.authenticate(ctx, userId, password); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("server", cfg.Server), errs.WithContext("user_id", userId))
	}
	return cfg, nil
}

func (cfg *Mastodon) register(ctx context.Context) error {
	app, err := mstdn.RegisterApp(ctx, &mstdn.AppConfig{
		Server:     cfg.Server,
		ClientName: cfg.AppName(),
		Scopes:     cfg.Scopes(),
		Website:    cfg.Registory(),
	})
	if err != nil {
		return errs.Wrap(err, errs.WithContext("server", cfg.Server))
	}
	cfg.Logger().Info().Interface("application", app).Str("server", cfg.Server).Msg("register application")
	cfg.ClientID = app.ClientID
	cfg.ClientSecret = app.ClientSecret
	return nil
}

func (cfg *Mastodon) authenticate(ctx context.Context, userId, password string) error {
	client := mstdn.NewClient(&mstdn.Config{
		Server:       cfg.Server,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	})
	cfg.Logger().Trace().Str("server", cfg.Server).Msg("start authntication")
	if err := client.Authenticate(ctx, userId, password); err != nil {
		return errs.Wrap(err, errs.WithContext("server", cfg.Server))
	}
	cfg.Logger().Trace().Str("server", cfg.Server).Msg("complete authntication")
	cfg.AccessToken = client.Config.AccessToken
	cfg.client = client
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
