package facade

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			apd, err := gopts.getAPOD(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			utcFlag, err := cmd.Flags().GetBool("utc")
			if err != nil {
				return debugPrint(ui, err)
			}
			dateStr, err := cmd.Flags().GetString("date")
			if err != nil {
				return debugPrint(ui, err)
			}
			date, err := values.DateFrom(dateStr, utcFlag)
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
			forceFlag, err := cmd.Flags().GetBool("force")
			if err != nil {
				return debugPrint(ui, err)
			}

			// lookup APOD data
			res, err := apd.LookupWithoutCache(cmd.Context(), date, utcFlag, forceFlag)
			if err != nil {
				apd.Logger().Error("error in apod.Lookup", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}

			// get image file
			fname, err := res.ImageFile(cmd.Context(), gopts.CacheDir)
			if err != nil && !errs.Is(err, ecode.ErrNoAPODImage) {
				return debugPrint(ui, err)
			}
			var imgs []string
			if len(fname) > 0 {
				defer os.Remove(fname)
				imgs = []string{fname}
			}

			// make message
			msg := apod.MakeMessage(res)

			var lastErrs []error

			// post to Bluesky
			if bskyFlag {
				wp, err := gopts.getWebpage(cmd.Context())
				if err != nil {
					return debugPrint(ui, err)
				}
				if bsky, err := gopts.getBluesky(wp); err != nil {
					apd.Logger().Info("no Bluesky configuration", zap.Object("error", zapobject.New(err)))
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
					apd.Logger().Info("no Mastodon configuration", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else if resText, err := mstdn.PostMessage(cmd.Context(), &mastodon.Message{
					Msg:        msg,
					ImageFiles: imgs,
				}); err != nil {
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
	apodPostCmd.Flags().BoolP("bluesky", "b", false, "Post to bluesky")
	apodPostCmd.Flags().BoolP("mastodon", "m", false, "Post to Mastodon")
	apodPostCmd.Flags().BoolP("force", "", false, "Force getting APOD data from cache")

	return apodPostCmd
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
