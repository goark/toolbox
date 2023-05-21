package bluesky

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"go.uber.org/zap"
)

// Profile method returns actor's profile information.
func (cfg *Bluesky) Profile(ctx context.Context, actor string) (*bsky.ActorDefs_ProfileViewDetailed, error) {
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
	cfg.Logger().Info("start getting profile", zap.String("actor", actor))
	profile, err := bsky.ActorGetProfile(ctx, cfg.client, actor)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("actor", actor))
	}
	cfg.Logger().Info("complete getting profile", zap.Any("profile", profile))
	return profile, nil
}

// ShowProfile method outouts actor's profile information to io.Wtiter.
func (cfg *Bluesky) ShowProfile(ctx context.Context, actor string, jsonFlag bool, w io.Writer) error {
	prof, err := cfg.Profile(ctx, actor)
	if err != nil {
		return errs.Wrap(err)
	}
	if jsonFlag {
		if err := json.NewEncoder(w).Encode(prof); err != nil {
			return errs.Wrap(err, errs.WithContext("actor", actor))
		}
	} else {
		fmt.Fprintf(w, " Handle Name: %s\n", prof.Handle)
		fmt.Fprintf(w, "         DID: %s\n", prof.Did)
		if prof.DisplayName != nil {
			fmt.Fprintf(w, "Display Name: %s\n", *prof.DisplayName)
		}
		if prof.IndexedAt != nil {
			fmt.Fprintf(w, "    Index at: %s\n", *prof.IndexedAt)
		}
		if prof.PostsCount != nil {
			fmt.Fprintf(w, "       Posts: %d\n", *prof.PostsCount)
		}
		if prof.FollowsCount != nil {
			fmt.Fprintf(w, "     Follows: %d\n", *prof.FollowsCount)
		}
		if prof.FollowersCount != nil {
			fmt.Fprintf(w, "   Followers: %d\n", *prof.FollowersCount)
		}
		if prof.Description != nil {
			fmt.Fprintf(w, "\n%s\n", *prof.Description)
		}
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
