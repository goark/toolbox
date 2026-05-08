package facade

import (
	"os"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/apod"
	"github.com/goark/toolbox/bluesky"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/mastodon"
	"github.com/goark/toolbox/values"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newAPODPostCmd returns cobra.Command instance for show sub-command
func newAPODPostCmd(ui *rwi.RWI) *cobra.Command {
	apodPostCmd := &cobra.Command{
		Use:     "post",
		Aliases: []string{"pst", "p"},
		Short:   "Post APOD data to TL",
		Long:    "Post Astronomy Picture of the Day data to time lines.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Global options
			gopts, oerr := getGlobalOptions()
			if oerr != nil {
				err = debugPrint(ui, oerr)
				return
			}
			apd, gerr := gopts.getAPOD(cmd.Context())
			if gerr != nil {
				err = debugPrint(ui, gerr)
				return
			}
			// local options
			utcFlag, ferr := cmd.Flags().GetBool("utc")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
			dateStr, ferr := cmd.Flags().GetString("date")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}
			date, verr := values.DateFrom(dateStr, utcFlag)
			if verr != nil {
				err = debugPrint(ui, verr)
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
			forceFlag, ferr := cmd.Flags().GetBool("force")
			if ferr != nil {
				err = debugPrint(ui, ferr)
				return
			}

			// lookup APOD data
			res, aerr := apd.LookupWithoutCache(cmd.Context(), date, utcFlag, forceFlag)
			if aerr != nil {
				apd.Logger().Error("error in apod.Lookup", zap.Object("error", zapobject.New(aerr)))
				err = debugPrint(ui, aerr)
				return
			}

			// get image file
			fname, rerr := res.ImageFile(cmd.Context(), gopts.CacheDir)
			if rerr != nil && !errs.Is(rerr, ecode.ErrNoAPODImage) {
				err = debugPrint(ui, rerr)
				return
			}
			var imgs []string
			if len(fname) > 0 {
				defer func() {
					err = errs.Join(err, os.Remove(fname))
				}()
				imgs = []string{fname}
			}

			// make message
			msg := apod.MakeMessage(res)

			var lastErrs []error

			// post to Bluesky
			if bskyFlag {
				wp, gerr := gopts.getWebpage(cmd.Context())
				if gerr != nil {
					err = debugPrint(ui, gerr)
					return
				}
				if bsky, gerr := gopts.getBluesky(wp); gerr != nil {
					apd.Logger().Info("no Bluesky configuration", zap.Object("error", zapobject.New(gerr)))
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
					apd.Logger().Info("no Mastodon configuration", zap.Object("error", zapobject.New(gerr)))
					lastErrs = append(lastErrs, gerr)
				} else if resText, merr := mstdn.PostMessage(cmd.Context(), &mastodon.Message{
					Msg:        msg,
					ImageFiles: imgs,
				}); merr != nil {
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
			return
		},
	}
	apodPostCmd.Flags().BoolP("bluesky", "b", false, "Post to bluesky")
	apodPostCmd.Flags().BoolP("mastodon", "m", false, "Post to Mastodon")
	apodPostCmd.Flags().BoolP("force", "", false, "Force getting APOD data from cache")

	return apodPostCmd
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
