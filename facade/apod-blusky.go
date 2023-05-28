package facade

import (
	"os"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/bluesky"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/values"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newAPODLookupCmd returns cobra.Command instance for show sub-command
func newAPODBlueskyCmd(ui *rwi.RWI) *cobra.Command {
	apodBlueskyCmd := &cobra.Command{
		Use:     "bluesky",
		Aliases: []string{"bsky", "bs"},
		Short:   "Post APOD data to Bluesky",
		Long:    "Post Astronomy Picture of the Day data to Bluesky.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			apd, err := getAPOD()
			if err != nil {
				return debugPrint(ui, err)
			}
			bsky, err := getBluesky()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			dateStr, err := cmd.Flags().GetString("date")
			if err != nil {
				return debugPrint(ui, err)
			}
			date, err := values.DateFrom(dateStr)
			if err != nil {
				return debugPrint(ui, err)
			}

			// lookup APOD data
			res, err := apd.Lookup(cmd.Context(), date)
			if err != nil {
				apd.Logger().Error("error in apod.Lookup", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}

			// make message
			credit := ""
			if len(res.Copyright) > 0 {
				credit = "\nImage Credit: " + res.Copyright
			}
			msg := strings.Join([]string{
				"#apod",
				res.Title + credit,
				res.WebPage(),
			}, "\n")

			// image file
			fname, err := res.ImageFile(cmd.Context(), bsky.BaseDir())
			if err != nil && !errs.Is(err, ecode.ErrNoAPODImage) {
				return debugPrint(ui, err)
			}
			var imgs []string
			if len(fname) > 0 {
				defer os.Remove(fname)
				imgs = []string{fname}
			}

			// post to bluesky
			resText, err := bsky.PostMessage(cmd.Context(), &bluesky.Message{Msg: msg, ImageFiles: imgs})
			if err != nil {
				bsky.Logger().Error("error in bluesky.PostMessage", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			return debugPrint(ui, ui.Outputln(resText))
		},
	}
	apodBlueskyCmd.Flags().StringP("date", "d", "", "Date for APOD data (YYYY-MM-DD)")

	return apodBlueskyCmd
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
