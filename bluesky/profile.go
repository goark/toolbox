package bluesky

import (
	"context"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
)

func (cfg *Bluesky) ShowProfile(ctx context.Context, actor string) (*bsky.ActorDefs_ProfileViewDetailed, error) {
	if cfg == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer, errs.WithContext("actor", actor))
	}

	// create/refresh session
	if cfg.client == nil {
		if err := cfg.CreateSession(ctx); err != nil {
			return nil, errs.Wrap(err, errs.WithContext("actor", actor))
		}
	}

	// get profile
	if len(actor) == 0 {
		actor = cfg.Handle
	}
	profile, err := bsky.ActorGetProfile(ctx, cfg.client, actor)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("actor", actor))
	}
	cfg.Logger().Info().Interface("profile", profile).Send()
	return profile, nil
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
