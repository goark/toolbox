package bluesky

import (
	"context"
	"net/url"

	"github.com/goark/errs"
	"github.com/ipfs/go-log/v2"
)

// Register functions registers application to mastodon server.
func Register(ctx context.Context, server, handle, password, baseDir string, logger *log.ZapEventLogger) (*Bluesky, error) {
	if len(server) == 0 {
		server = DefaltHostName
	} else if u, err := url.Parse(server); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("server", server))
	} else if len(u.Hostname()) > 0 {
		server = u.Hostname()
	}
	cfg := &Bluesky{
		Host:     "https://" + server,
		Handle:   handle,
		Password: password,
		baseDir:  baseDir,
		logger:   logger,
	}

	// create session
	if err := cfg.CreateSession(ctx); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("host", server), errs.WithContext("handle", handle))
	}
	return cfg, nil
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
