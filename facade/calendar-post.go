package facade

import (
	"bytes"
	"errors"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/bluesky"
	"github.com/goark/toolbox/mastodon"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newCalendarPostCmd returns cobra.Command instance for show sub-command
func newCalendarPostCmd(ui *rwi.RWI) *cobra.Command {
	calendarPostCmd := &cobra.Command{
		Use:     "post",
		Aliases: []string{"pst", "p"},
		Short:   "Post astronomical calendar data to TL",
		Long:    "Post astronomical calendar data to time lines.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			// make temporary directory
			if err := gopts.TempDir.MakeDir(); err != nil {
				return debugPrint(ui, err)
			}
			defer func() { _ = gopts.TempDir.CleanUp() }()
			// local options
			ccfg, err := getCalendarConfig(cmd, gopts)
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

			// lookup calendar data
			b := &bytes.Buffer{}
			if err := ccfg.OutputEvent(b); err != nil {
				gopts.Logger.Error("error in calendar.Config.OutputEvent", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			msg := b.String()

			var lastErrs []error

			// post to Bluesky
			if bskyFlag {
				wp, err := gopts.getWebpage(cmd.Context())
				if err != nil {
					return debugPrint(ui, err)
				}
				if bsky, err := gopts.getBluesky(wp); err != nil {
					gopts.Logger.Info("no Bluesky configuration", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else if resText, err := bsky.PostMessage(cmd.Context(), &bluesky.Message{Msg: msg, ImageFiles: nil}); err != nil {
					bsky.Logger().Error("error in bluesky.PostMessage", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else {
					_ = ui.Outputln("post to Bluesky:", resText)
				}
			}
			// post to Mastodon
			if mastodonFlag {
				if mstdn, err := gopts.getMastodon(); err != nil {
					gopts.Logger.Info("no Mastodon configuration", zap.Object("error", zapobject.New(err)))
					lastErrs = append(lastErrs, err)
				} else if resText, err := mstdn.PostMessage(cmd.Context(), &mastodon.Message{
					Msg:        msg,
					ImageFiles: nil,
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
	calendarPostCmd.Flags().BoolP("bluesky", "b", false, "Post to bluesky")
	calendarPostCmd.Flags().BoolP("mastodon", "m", false, "Post to Mastodon")

	return calendarPostCmd
}

/* Copyright 2024 Spiegel
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
