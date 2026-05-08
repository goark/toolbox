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

// newBookmarkPostCmd returns cobra.Command instance for show sub-command
func newBookmarkPostCmd(ui *rwi.RWI) *cobra.Command {
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
			urlStr, ferr := cmd.Flags().GetString("url")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
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

			// lookup Web page data
			page, cerr := cfg.Lookup(cmd.Context(), urlStr)
			if cerr != nil {
				gopts.Logger.Desugar().Error("error in bookmark.Lookup", zap.Object("error", zapobject.New(cerr)))
				err = debugPrint(ui, cerr)
				return
			}

			// get image file
			gopts.Logger.Desugar().Debug("start posting web page info", zap.Any("info", page))
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

			var lastErrs []error

			// post to Bluesky
			if bskyFlag {
				if bsky, gerr := gopts.getBluesky(cfg); gerr != nil {
					cfg.Logger().Info("no Bluesky configuration", zap.Object("error", zapobject.New(gerr)))
					lastErrs = append(lastErrs, gerr)
				} else if resText, berr := bsky.PostMessage(cmd.Context(), &bluesky.Message{Msg: msg, ImageFiles: imgs}); berr != nil {
					bsky.Logger().Error("error in bluesky.PostMessage", zap.Object("error", zapobject.New(berr)))
					lastErrs = append(lastErrs, berr)
				} else {
					_ = ui.Outputln("post to Bluesky:", resText)
				}
			}
			// post to Mastodon
			if mastodonFlag {
				if mstdn, gerr := gopts.getMastodon(); gerr != nil {
					cfg.Logger().Info("no Mastodon configuration", zap.Object("error", zapobject.New(gerr)))
					lastErrs = append(lastErrs, gerr)
				} else if resText, merr := mstdn.PostMessage(cmd.Context(), &mastodon.Message{Msg: msg, ImageFiles: imgs}); merr != nil {
					mstdn.Logger().Error("error in mastodon.PostMessage", zap.Object("error", zapobject.New(merr)))
					lastErrs = append(lastErrs, merr)
				} else {
					_ = ui.Outputln("post to Mastodon:", resText)
				}
			}

			if len(lastErrs) > 0 {
				err = debugPrint(ui, errs.Wrap(errs.Join(lastErrs...)))
				return
			}
			gopts.Logger.Desugar().Debug("end posting web page info", zap.Any("page", page))

			if saveFlag {
				list, cerr := cfg.StopPool()
				if cerr != nil {
					err = debugPrint(ui, cerr)
					return
				}
				if serr := cfg.Save(cmd.Context(), list); serr != nil {
					gopts.Logger.Desugar().Error("error in webpage.Lookup", zap.Object("error", zapobject.New(serr)))
					err = debugPrint(ui, serr)
					return
				}
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
