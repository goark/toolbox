package facade

import (
	"errors"
	"os"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/bluesky"
	"github.com/goark/toolbox/bookmark"
	"github.com/goark/toolbox/mastodon"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newBookmarkPostCmd returns cobra.Command instance for show sub-command
func newBookmarkPostCmd(ui *rwi.RWI) *cobra.Command {
	bookmarkPostCmd := &cobra.Command{
		Use:     "post",
		Aliases: []string{"pst", "p"},
		Short:   "Post Web page's information to Bluesky",
		Long:    "Post Web page's information to Bluesky.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			cfg, err := gopts.getBookmark()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			urlStr, err := cmd.Flags().GetString("url")
			if err != nil {
				return debugPrint(ui, err)
			}
			saveFlag, err := cmd.Flags().GetBool("save")
			if err != nil {
				return debugPrint(ui, err)
			}
			bskyFlag, err := cmd.Flags().GetBool("bluesky")
			if err != nil {
				return debugPrint(ui, err)
			}
			mastodonFlag, err := cmd.Flags().GetBool("mastodon")
			if err != nil {
				return debugPrint(ui, err)
			}
			withImage, err := cmd.Flags().GetBool("with-image")
			if err != nil {
				return debugPrint(ui, err)
			}
			pmsg, err := cmd.Flags().GetString("prefix-message")
			if err != nil {
				return debugPrint(ui, err)
			}

			// lookup Web page data
			info, err := cfg.Lookup(cmd.Context(), urlStr, saveFlag)
			if err != nil {
				gopts.Logger.Desugar().Error("error in bookmark.Lookup", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}

			// get image file
			var imgs []string
			if withImage && len(info.ImageURL) > 0 {
				fname, err := info.ImageFile(cmd.Context(), gopts.CacheDir)
				if err != nil {
					return debugPrint(ui, err)
				}
				if len(fname) > 0 {
					defer os.Remove(fname)
					imgs = []string{fname}
				}
			}

			// make message
			msg := bookmark.MakeMessage(info, strings.TrimSpace(pmsg))

			var lastErrs []error

			// post to Bluesky
			if bskyFlag {
				if bsky, err := gopts.getBluesky(); err != nil {
					cfg.Logger().Info("no Bluesky configuration", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else if resText, err := bsky.PostMessage(cmd.Context(), &bluesky.Message{Msg: msg, ImageFiles: imgs}); err != nil {
					bsky.Logger().Error("error in bluesky.PostMessage", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else {
					_ = ui.Outputln("post to Bluesky:", resText)
				}
			}
			// post to Mastodon
			if mastodonFlag {
				if mstdn, err := gopts.getMastodon(); err != nil {
					cfg.Logger().Info("no Mastodon configuration", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else if resText, err := mstdn.PostMessage(cmd.Context(), &mastodon.Message{Msg: msg, ImageFiles: imgs}); err != nil {
					mstdn.Logger().Error("error in mastodon.PostMessage", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else {
					_ = ui.Outputln("post to Mastodon:", resText)
				}
			}

			if len(lastErrs) > 0 {
				return debugPrint(ui, errs.Wrap(errors.Join(lastErrs...)))
			}
			return nil
		},
	}
	bookmarkPostCmd.Flags().BoolP("bluesky", "b", false, "Post to bluesky")
	bookmarkPostCmd.Flags().BoolP("mastodon", "m", false, "Post to Mastodon")
	bookmarkPostCmd.Flags().BoolP("with-image", "", false, "Post with image")
	bookmarkPostCmd.Flags().StringP("prefix-message", "p", "", "Message")

	return bookmarkPostCmd
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
