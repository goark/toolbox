package bluesky

import (
	"context"
	"net/url"
	"strings"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"go.uber.org/zap"
)

func (cfg *Bluesky) getRecord(ctx context.Context, uri string) (*atproto.RepoGetRecord_Output, error) {
	var did, collection, recordKey string
	if strings.HasPrefix(uri, "at://") {
		// record URI: at://{did}/{collection}/{recordKey}
		parts := strings.Split(strings.TrimPrefix(uri, "at://"), "/")
		if len(parts) != 3 {
			return nil, errs.Wrap(ecode.ErrInvalidBlueskyRecordURI, errs.WithContext("uri", uri), errs.WithContext("detail", "error in split of URL(at://~)"))
		}
		did = parts[0]
		collection = parts[1]
		recordKey = parts[2]
	} else {
		// maybe http(s)://{host}/profile/{handle}/post/{recordKey}
		u, err := url.Parse(uri)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("uri", uri))
		}
		parts := strings.Split(u.Path, "/")
		if len(parts) < 5 {
			return nil, errs.Wrap(ecode.ErrInvalidBlueskyRecordURI, errs.WithContext("uri", uri), errs.WithContext("parts", parts))
		}
		if parts[1] != "profile" {
			return nil, errs.Wrap(ecode.ErrInvalidBlueskyRecordURI, errs.WithContext("uri", uri), errs.WithContext("detail", "parse error (profile)"))
		}
		p, err := cfg.Profile(ctx, parts[2])
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("uri", uri), errs.WithContext("handle", parts[2]))
		}
		did = p.Did
		if parts[3] != "post" {
			return nil, errs.Wrap(ecode.ErrInvalidBlueskyRecordURI, errs.WithContext("uri", uri), errs.WithContext("detail", "parse error (post)"))
		}
		recordKey = parts[4]
		collection = "app.bsky.feed.post"
	}
	cfg.Logger().Debug("start getting record", zap.String("collection", collection), zap.String("record_key", recordKey))
	record, err := atproto.RepoGetRecord(context.TODO(), cfg.client, "", collection, did, recordKey)
	if err != nil {
		return nil, errs.Wrap(
			err,
			errs.WithContext("uri", uri),
			errs.WithContext("did", did),
			errs.WithContext("collection", collection),
			errs.WithContext("record_key", recordKey),
		)
	}
	cfg.Logger().Info("complete getting record", zap.Any("record", record))
	return record, nil
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
