package facade

import (
	"strings"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/mastodon"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newMastodonPostCmd returns cobra.Command instance for show sub-command
func newMastodonPostCmd(ui *rwi.RWI) *cobra.Command {
	mastodonPostCmd := &cobra.Command{
		Use:     "post",
		Aliases: []string{"pst", "p", "toot", "tt", "t"},
		Short:   "Post message to Mastodon",
		Long:    "Post message to Mastodon.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			mstdn, err := gopts.getMastodon()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			images, err := cmd.Flags().GetStringSlice("image-file")
			if err != nil {
				return debugPrint(ui, err)
			}
			visStr, err := cmd.Flags().GetString("visibility")
			if err != nil {
				return debugPrint(ui, err)
			}
			visibility := mastodon.GetVisibilityFrom(visStr)
			if visibility == mastodon.VisibilityUnknown {
				return debugPrint(ui, errs.New("invlid visibility", errs.WithContext("visibility", visStr)))
			}
			spoilerText, err := cmd.Flags().GetString("spoiler-text")
			if err != nil {
				return debugPrint(ui, err)
			}
			msg, err := cmd.Flags().GetString("text")
			if err != nil {
				return debugPrint(ui, err)
			}
			pipeFlag, err := cmd.Flags().GetBool("pipe")
			if err != nil {
				return debugPrint(ui, err)
			}
			editFlag, err := cmd.Flags().GetBool("edit")
			if err != nil {
				return debugPrint(ui, err)
			}
			if pipeFlag {
				msg, err = inputFromPipe(ui)
				if err != nil {
					return debugPrint(ui, err)
				}
			} else if editFlag {
				msg, err = editMessage(cmd.Context(), ui.Writer())
				if err != nil {
					return debugPrint(ui, err)
				}
			}
			msg = strings.TrimSpace(msg)

			// post message
			resText, err := mstdn.PostMessage(cmd.Context(), &mastodon.Message{
				Msg:         msg,
				SpoilerText: spoilerText,
				Visibility:  visibility.String(),
				ImageFiles:  images,
			})
			if err != nil {
				mstdn.Logger().Error("error in mastodon.PostMessage", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			return debugPrint(ui, ui.Outputln(resText))
		},
	}
	mastodonPostCmd.Flags().StringP("text", "t", "", "Text message")
	mastodonPostCmd.Flags().BoolP("pipe", "", false, "Input from standard-input")
	mastodonPostCmd.Flags().BoolP("edit", "", false, "Edit message")
	mastodonPostCmd.MarkFlagsMutuallyExclusive("text", "pipe", "edit")
	mastodonPostCmd.Flags().StringSliceP("image-file", "i", nil, "Image file")
	mastodonPostCmd.Flags().StringP("visibility", "v", mastodon.DefaultVisibility().String(), "Visibility ["+strings.Join(mastodon.VisibilityList(), "|")+"]")
	mastodonPostCmd.Flags().StringP("spoiler-text", "s", "", "Spoiler text")

	return mastodonPostCmd
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
