package mastodon

import (
	"context"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	mstdn "github.com/mattn/go-mastodon"
)

type Visibility int

const (
	VisibilityUnknown = iota
	VisibilityPublic
	VisibilityUnlisted
	VisibilityFollowersOnly
	VisibilityDirectMessage
)

var (
	visibilityMap = map[Visibility]string{
		VisibilityPublic:        mstdn.VisibilityPublic,
		VisibilityUnlisted:      mstdn.VisibilityUnlisted,
		VisibilityFollowersOnly: mstdn.VisibilityFollowersOnly,
		VisibilityDirectMessage: mstdn.VisibilityDirectMessage,
	}
	visibilityList = []string{
		visibilityMap[VisibilityPublic],
		visibilityMap[VisibilityUnlisted],
		visibilityMap[VisibilityFollowersOnly],
		visibilityMap[VisibilityDirectMessage],
	}
)

// VisibilityList function returns list of Visibility strings.
func VisibilityList() []string {
	return visibilityList
}

// DefaultVisibility returns default Visibility.
func DefaultVisibility() Visibility {
	return VisibilityPublic
}

// GetVisibilityFrom function returns Visibility from string.
func GetVisibilityFrom(s string) Visibility {
	if len(s) == 0 {
		return DefaultVisibility()
	}
	for k, v := range visibilityMap {
		if strings.EqualFold(v, s) {
			return k
		}
	}
	return VisibilityUnknown
}

func (v Visibility) String() string {
	if s, ok := visibilityMap[v]; ok {
		return s
	}
	return ""
}

// Message is information of post message.
type Message struct {
	Msg         string
	SpoilerText string
	Visibility  string
	ImageFiles  []string
}

// PostMessage method posts message and image files to Mastodon.
func (cfg *Mastodon) PostMessage(ctx context.Context, msg *Message) (string, error) {
	if cfg == nil || cfg.client == nil {
		return "", errs.Wrap(ecode.ErrNullPointer)
	}

	// upload images
	images, err := cfg.uploadImages(ctx, msg.ImageFiles)
	if err != nil {
		return "", errs.Wrap(err)
	}

	// make toot
	toot := &mstdn.Toot{
		Status:      msg.Msg,
		Visibility:  msg.Visibility,
		SpoilerText: msg.SpoilerText,
		MediaIDs:    images,
	}
	if len(toot.SpoilerText) > 0 {
		toot.Sensitive = true
	}
	cfg.Logger().Debug().Interface("toot", toot).Send()

	// post toot
	cfg.Logger().Debug().Msg("start posting message")
	stat, err := cfg.client.PostStatus(ctx, toot)
	if err != nil {
		return "", errs.Wrap(err)
	}
	cfg.Logger().Info().Interface("response_of_post", stat).Msg("complete posting message")
	return stat.URL, nil
}

func (cfg *Mastodon) uploadImages(ctx context.Context, paths []string) ([]mstdn.ID, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	list := make([]mstdn.ID, 0, len(paths))
	for _, path := range paths {
		cfg.Logger().Debug().Str("path", path).Msg("start uploading image file")
		attch, err := cfg.client.UploadMedia(ctx, path)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("path", path))
		}
		cfg.Logger().Info().Str("path", path).Interface("atach_info", attch).Msg("complete uploading image file")
		list = append(list, attch.ID)
	}
	return list, nil
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
