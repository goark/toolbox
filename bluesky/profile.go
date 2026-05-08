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
func (cfg *Bluesky) Profile(ctx context.Context, actor string) (profile *bsky.ActorDefs_ProfileViewDetailed, err error) {
	if cfg == nil {
		err = errs.Wrap(ecode.ErrNullPointer, errs.WithContext("actor", actor))
		return
	}

	// create/refresh session
	if cfg.client == nil {
		if cerr := cfg.CreateSession(ctx); cerr != nil {
			err = errs.Wrap(cerr, errs.WithContext("actor", actor))
			return
		}
	}

	// get profile
	if len(actor) == 0 {
		actor = cfg.Handle
	}
	cfg.Logger().Info("start getting profile", zap.String("actor", actor))
	profile, err = bsky.ActorGetProfile(ctx, cfg.client, actor)
	if err != nil {
		err = errs.Wrap(err, errs.WithContext("actor", actor))
		return
	}
	cfg.Logger().Info("complete getting profile", zap.Any("profile", profile))
	return
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
		if _, err := fmt.Fprintf(w, " Handle Name: %s\n", prof.Handle); err != nil {
			return errs.Wrap(err, errs.WithContext("actor", actor))
		}
		if _, err := fmt.Fprintf(w, "         DID: %s\n", prof.Did); err != nil {
			return errs.Wrap(err, errs.WithContext("actor", actor))
		}
		if prof.DisplayName != nil {
			if _, err := fmt.Fprintf(w, "Display Name: %s\n", *prof.DisplayName); err != nil {
				return errs.Wrap(err, errs.WithContext("actor", actor))
			}
		}
		if prof.IndexedAt != nil {
			if _, err := fmt.Fprintf(w, "    Index at: %s\n", *prof.IndexedAt); err != nil {
				return errs.Wrap(err, errs.WithContext("actor", actor))
			}
		}
		if prof.PostsCount != nil {
			if _, err := fmt.Fprintf(w, "       Posts: %d\n", *prof.PostsCount); err != nil {
				return errs.Wrap(err, errs.WithContext("actor", actor))
			}
		}
		if prof.FollowsCount != nil {
			if _, err := fmt.Fprintf(w, "     Follows: %d\n", *prof.FollowsCount); err != nil {
				return errs.Wrap(err, errs.WithContext("actor", actor))
			}
		}
		if prof.FollowersCount != nil {
			if _, err := fmt.Fprintf(w, "   Followers: %d\n", *prof.FollowersCount); err != nil {
				return errs.Wrap(err, errs.WithContext("actor", actor))
			}
		}
		if prof.Description != nil {
			if _, err := fmt.Fprintf(w, "\n%s\n", *prof.Description); err != nil {
				return errs.Wrap(err, errs.WithContext("actor", actor))
			}
		}
	}
	return nil
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
