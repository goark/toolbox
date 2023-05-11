package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	mstdn "github.com/mattn/go-mastodon"
)

// Profile method retuns account information of Mastodon user.
func (cfg *Mastodon) Profile(ctx context.Context) (*mstdn.Account, error) {
	if cfg == nil || cfg.client == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	account, err := cfg.client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return account, nil
}

// ShowProfile method outputs account information of Mastodon user to io.Writer.
func (cfg *Mastodon) ShowProfile(ctx context.Context, jsonFlag bool, w io.Writer) error {
	account, err := cfg.Profile(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	if jsonFlag {
		if err := json.NewEncoder(w).Encode(account); err != nil {
			return errs.Wrap(err)
		}
	} else {
		fmt.Fprintf(w, "      Username: %s\n", account.Username)
		fmt.Fprintf(w, "User ID (full): @%s@%s\n", account.Username, cfg.Servername())
		fmt.Fprintf(w, "           URL: %s\n", account.URL)
		fmt.Fprintf(w, "  Display name: %s\n", account.DisplayName)
		fmt.Fprintf(w, "    Created at: %v\n", account.CreatedAt)
		fmt.Fprintf(w, "         Posts: %d\n", account.StatusesCount)
		fmt.Fprintf(w, "       Follows: %d\n", account.FollowingCount)
		fmt.Fprintf(w, "     Followers: %d\n", account.FollowersCount)
		fmt.Fprintf(w, "\n%s\n", account.Note)
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
