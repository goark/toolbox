package bluesky

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/images"
	"github.com/goark/toolbox/webpage"
)

// Message is information of post message.
type Message struct {
	Msg        string
	ImageFiles []string
}

func (cfg *Bluesky) PostMessage(ctx context.Context, msg *Message) (string, error) {
	if cfg == nil {
		return "", errs.Wrap(ecode.ErrNullPointer, errs.WithContext("msg", msg))
	}
	if len(msg.Msg) == 0 {
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
		Text:      msg.Msg,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Reply:     nil,
	}

	// add links metadata
	for _, e := range getLinksFrom(msg.Msg) {
		post.Entities = append(post.Entities, &bsky.FeedPost_Entity{
			Index: &bsky.FeedPost_TextSlice{
				Start: e.start,
				End:   e.end,
			},
			Type:  "link",
			Value: e.text,
		})
		if post.Embed == nil {
			post.Embed = &bsky.FeedPost_Embed{}
		}
		if post.Embed.EmbedExternal == nil {
			// get information of web page
			if link, err := webpage.ReadPage(ctx, e.text); err != nil {
				cfg.Logger().Info().Interface("error", errs.Wrap(err)).Str("web_page", e.text).Msg("cannot read web page")
			} else {
				post.Embed.EmbedExternal = &bsky.EmbedExternal{
					External: &bsky.EmbedExternal_External{
						Description: link.Description,
						Title:       link.Title,
						Uri:         link.URL,
					},
				}
				cfg.Logger().Trace().Str("title", link.Title).Str("description", link.Description).Str("url", link.URL).Msg("web page info")
				// get attention image
				if len(link.ImageURL) > 0 {
					if res, err := cfg.getEmbedImage(ctx, link.ImageURL); err != nil {
						cfg.Logger().Info().Interface("error", errs.Wrap(err)).Str("image_url", link.ImageURL).Msg("cannot get embeded image")
					} else {
						post.Embed.EmbedExternal.External.Thumb = res.Blob
						cfg.Logger().Trace().Str("contentType", res.Blob.MimeType).Int64("size", res.Blob.Size).Str("url", link.ImageURL).Msg("embeded file")
					}
				}
			}
		}
	}

	// add mentions metadata
	for _, e := range getMentionsFrom(msg.Msg) {
		prof, err := cfg.Profile(ctx, e.text)
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

	// embeded images
	if len(msg.ImageFiles) > 0 {
		for _, fn := range msg.ImageFiles {
			var imgs []*bsky.EmbedImages_Image
			src, err := images.FetchFromFile(fn)
			if err != nil {
				return "", errs.Wrap(err, errs.WithContext("file", fn))
			}
			img, err := images.AjustImage(src)
			if err != nil {
				cfg.Logger().Error().Interface("error", errs.Wrap(err)).Str("file_name", fn).Msg("cannot ajust image")
				return "", errs.Wrap(err, errs.WithContext("file", fn))
			}
			cfg.Logger().Trace().Str("file_name", fn).Msg("start uploading image file")
			res, err := atproto.RepoUploadBlob(ctx, cfg.client, img)
			if err != nil {
				cfg.Logger().Error().Interface("error", errs.Wrap(err)).Str("file_name", fn).Msg("cannot upload image file")
				return "", errs.Wrap(err, errs.WithContext("fn", fn))
			}
			imgs = append(imgs, &bsky.EmbedImages_Image{
				Alt:   filepath.Base(fn),
				Image: res.Blob,
			})
			if post.Embed == nil {
				post.Embed = &bsky.FeedPost_Embed{}
			}
			cfg.Logger().Trace().Str("contentType", res.Blob.MimeType).Int64("size", res.Blob.Size).Str("file_name", fn).Msg("complete uploading image file")
			post.Embed.EmbedImages = &bsky.EmbedImages{Images: imgs}
		}
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

func (cfg *Bluesky) getEmbedImage(ctx context.Context, urlStr string) (*atproto.RepoUploadBlob_Output, error) {
	src, err := images.FetchFromURL(ctx, urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	img, err := images.AjustImage(src)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}

	res, err := atproto.RepoUploadBlob(ctx, cfg.client, img)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	return res, nil
}

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
