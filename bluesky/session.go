package bluesky

import (
	"context"
	"net/http"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
)

// CreateSession method makes XRPC instance and creates/refreshes session.
func (cfg *Bluesky) CreateSession(ctx context.Context) error {
	if cfg == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	if cfg.client != nil {
		cfg.Logger().Debug().Msg("exist session")
		return nil
	}
	client := &xrpc.Client{
		Client: defaultHttpClient(),
		Host:   cfg.Host,
		Auth:   &xrpc.AuthInfo{Handle: cfg.Handle},
	}
	auth, err := cfg.readAuth()
	if err == nil { // exist auth file
		// refresh session
		cfg.Logger().Debug().Msg("start refreshing session")
		client.Auth = auth
		client.Auth.AccessJwt = client.Auth.RefreshJwt
		refresh, err2 := atproto.ServerRefreshSession(ctx, client)
		if err2 != nil {
			err = err2 // --> create session
			cfg.Logger().Info().Msg("cannot refresh session.")
		} else {
			cfg.Logger().Debug().Msg("complete refreshing session")
			client.Auth.Did = refresh.Did
			client.Auth.AccessJwt = refresh.AccessJwt
			client.Auth.RefreshJwt = refresh.RefreshJwt
			if err := cfg.writeAuth(client.Auth); err != nil {
				return errs.Wrap(err)
			}
		}
	}
	if err != nil {
		cfg.Logger().Info().Interface("error", errs.Wrap(err)).Msg("no valid auth-file, start creating session")
		// create session
		auth, err := atproto.ServerCreateSession(ctx, client, &atproto.ServerCreateSession_Input{
			Identifier: client.Auth.Handle,
			Password:   cfg.Password,
		})
		if err != nil {
			return errs.Wrap(err)
		}
		cfg.Logger().Debug().Msg("complete creating session")
		client.Auth.Did = auth.Did
		client.Auth.AccessJwt = auth.AccessJwt
		client.Auth.RefreshJwt = auth.RefreshJwt
		if err := cfg.writeAuth(client.Auth); err != nil {
			return errs.Wrap(err)
		}
	}
	cfg.client = client
	return nil
}

// see github.com/bluesky-social/indigo/cmd/gosky/util package.
func defaultHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
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
