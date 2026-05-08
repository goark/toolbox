package facade

import (
	"os"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/bluesky"
	"github.com/goark/toolbox/mastodon"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newFeedPostCmd returns cobra.Command instance for show sub-command
func newFeedPostCmd(ui *rwi.RWI) *cobra.Command {
	bookmarkPostCmd := &cobra.Command{
		Use:     "post",
		Aliases: []string{"pst", "p"},
		Short:   "Post Web page's information to TL",
		Long:    "Post Web page's information to time lines.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Global options
			gopts, gerr := getGlobalOptions()
			if gerr != nil {
				err = debugPrint(ui, gerr)
				return
			}
			cfg, gerr := gopts.getWebpage(cmd.Context())
			if gerr != nil {
				err = debugPrint(ui, gerr)
				return
			}
			// local options
			saveFlag, ferr := cmd.Flags().GetBool("save")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
			bskyFlag, ferr := cmd.Flags().GetBool("bluesky")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
			mastodonFlag, ferr := cmd.Flags().GetBool("mastodon")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
			withImage, ferr := cmd.Flags().GetBool("with-image")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
			pmsg, ferr := cmd.Flags().GetString("prefix-text")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}

			// lookup feed
			list, ferr := getFeedAll(cmd, cfg)
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
			if saveFlag && len(list) > 0 {
				if serr := cfg.Save(cmd.Context(), list); serr != nil {
					err = debugPrint(ui, serr)
					return
				}
			}

			// post feed data
			var lastErrs []error
			var bsky *bluesky.Bluesky
			var mstdn *mastodon.Mastodon
			for _, page := range list {
				gopts.Logger.Desugar().Debug("start posting web page info", zap.Any("info", page))
				// get image file
				var imgs []string
				if withImage && len(page.ImageURL) > 0 {
					fname, ferr := page.ImageFile(cmd.Context(), gopts.CacheDir)
					if ferr != nil {
						err = debugPrint(ui, ferr)
						return
					}
					if len(fname) > 0 {
						gopts.Logger.Desugar().Debug("downloaded image file", zap.String("url", page.ImageURL), zap.String("local", fname))
						defer func() {
							err = errs.Join(err, os.Remove(fname))
						}()
						imgs = []string{fname}
					}
				}
				// make message
				msg := page.MakeMessage(strings.TrimSpace(pmsg))
				// post to Bluesky
				if bskyFlag {
					if bsky == nil {
						bsky, gerr = gopts.getBluesky(cfg)
						if gerr != nil {
							cfg.Logger().Info("no Bluesky configuration", zap.Object("error", zapobject.New(gerr)))
							err = debugPrint(ui, gerr)
							return
						}
					}
					if resText, berr := bsky.PostMessage(cmd.Context(), &bluesky.Message{Msg: msg, ImageFiles: imgs}); berr != nil {
						bsky.Logger().Error("error in bluesky.PostMessage", zap.Object("error", zapobject.New(berr)))
						lastErrs = append(lastErrs, berr)
					} else {
						_ = ui.Outputln("post to Bluesky:", resText)
					}
				}
				// post to Mastodon
				if mastodonFlag {
					if mstdn == nil {
						mstdn, gerr = gopts.getMastodon()
						if gerr != nil {
							cfg.Logger().Info("no Mastodon configuration", zap.Object("error", zapobject.New(gerr)))
							err = debugPrint(ui, gerr)
							return
						}
					}
					if resText, merr := mstdn.PostMessage(cmd.Context(), &mastodon.Message{Msg: msg, ImageFiles: imgs}); merr != nil {
						mstdn.Logger().Error("error in mastodon.PostMessage", zap.Object("error", zapobject.New(merr)))
						lastErrs = append(lastErrs, merr)
					} else {
						_ = ui.Outputln("post to Mastodon:", resText)
					}
				}
				gopts.Logger.Desugar().Debug("end posting web page info", zap.Any("info", page))
			}

			if len(lastErrs) > 0 {
				err = debugPrint(ui, errs.Wrap(errs.Join(lastErrs...)))
				return
			}
			return
		},
	}
	bookmarkPostCmd.Flags().BoolP("bluesky", "b", false, "Post to bluesky")
	bookmarkPostCmd.Flags().BoolP("mastodon", "m", false, "Post to Mastodon")
	bookmarkPostCmd.Flags().BoolP("with-image", "", false, "Post with image")
	bookmarkPostCmd.Flags().StringP("prefix-text", "t", "", "prefix text message")

	return bookmarkPostCmd
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
