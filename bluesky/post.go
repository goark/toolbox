package bluesky

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
)

var (
	urlRegexp     = regexp.MustCompile(`https?://[-A-Za-z0-9+&@#\/%?=~_|!:,.;\(\)]+`)
	mentionRegexp = regexp.MustCompile(`@[a-zA-Z0-9.]+`)
)

type entry struct {
	start int64
	end   int64
	text  string
}

func getLinksFrom(msg string) []entry {
	var result []entry
	matches := urlRegexp.FindAllStringSubmatchIndex(msg, -1)
	for _, m := range matches {
		result = append(result, entry{
			text:  msg[m[0]:m[1]],
			start: int64(len([]rune(msg[0:m[0]]))),
			end:   int64(len([]rune(msg[0:m[1]])))},
		)
	}
	return result
}

func getMentionsFrom(msg string) []entry {
	var result []entry
	matches := mentionRegexp.FindAllStringSubmatchIndex(msg, -1)
	for _, m := range matches {
		result = append(result, entry{
			text:  strings.TrimPrefix(msg[m[0]:m[1]], "@"),
			start: int64(len([]rune(msg[0:m[0]]))),
			end:   int64(len([]rune(msg[0:m[1]])))},
		)
	}
	return result
}

func (cfg *Bluesky) PostMessage(ctx context.Context, msg string) (string, error) {
	if cfg == nil {
		return "", errs.Wrap(ecode.ErrNullPointer, errs.WithContext("msg", msg))
	}
	if len(msg) == 0 {
		return "", errs.Wrap(ecode.ErrNoContent, errs.WithContext("msg", msg))
	}

	// create/refresh session
	if cfg.client == nil {
		if err := cfg.CreateSession(ctx); err != nil {
			return "", errs.Wrap(err, errs.WithContext("msg", msg))
		}
	}

	// make post data (not support reply and quote)
	post := &bsky.FeedPost{
		Text:      msg,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Reply:     nil,
	}

	for _, e := range getLinksFrom(msg) {
		post.Entities = append(post.Entities, &bsky.FeedPost_Entity{
			Index: &bsky.FeedPost_TextSlice{
				Start: e.start,
				End:   e.end,
			},
			Type:  "link",
			Value: e.text,
		})
	}

	// add mentions metadata
	for _, e := range getMentionsFrom(msg) {
		prof, err := cfg.ShowProfile(ctx, e.text)
		if err != nil {
			return "", errs.Wrap(err, errs.WithContext("msg", msg))
		}
		post.Entities = append(post.Entities, &bsky.FeedPost_Entity{
			Index: &bsky.FeedPost_TextSlice{
				Start: e.start,
				End:   e.end,
			},
			Type:  "mention",
			Value: prof.Did,
		})
	}

	// pos message
	cfg.Logger().Debug().Msg("start posting message")
	resp, err := atproto.RepoCreateRecord(ctx, cfg.client, &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       cfg.client.Auth.Did,
		Record: &util.LexiconTypeDecoder{
			Val: post,
		},
	})
	if err != nil {
		return "", errs.Wrap(err, errs.WithContext("msg", msg))
	}
	cfg.Logger().Info().Interface("response_of_post", resp).Msg("complete posting message")

	return resp.Uri, nil
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
