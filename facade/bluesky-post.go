package facade

import (
	"strings"

	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/bluesky"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newBlueskyPostCmd returns cobra.Command instance for show sub-command
func newBlueskyPostCmd(ui *rwi.RWI) *cobra.Command {
	blueskyPostCmd := &cobra.Command{
		Use:     "post",
		Aliases: []string{"pst", "p"},
		Short:   "Post message to Bluesky",
		Long:    "Post message to Bluesky.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			bsky, err := getBluesky()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			images, err := cmd.Flags().GetStringSlice("image-file")
			if err != nil {
				return debugPrint(ui, err)
			}
			replyTo, err := cmd.Flags().GetString("reply-to")
			if err != nil {
				return debugPrint(ui, err)
			}
			replyTo = strings.TrimSpace(replyTo)
			msg, err := cmd.Flags().GetString("message")
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
			resText, err := bsky.PostMessage(cmd.Context(), &bluesky.Message{Msg: msg, ImageFiles: images, ReplryTo: replyTo})
			if err != nil {
				bsky.Logger().Error("error in bluesky.PostMessage", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			return debugPrint(ui, ui.Outputln(resText))
		},
	}
	blueskyPostCmd.Flags().StringP("message", "m", "", "Message")
	blueskyPostCmd.Flags().BoolP("pipe", "", false, "Input from standard-input")
	blueskyPostCmd.Flags().BoolP("edit", "", false, "Edit message")
	blueskyPostCmd.MarkFlagsMutuallyExclusive("message", "pipe", "edit")
	blueskyPostCmd.Flags().StringSliceP("image-file", "i", nil, "Image file")
	blueskyPostCmd.Flags().StringP("reply-to", "r", "", "Replry URI")
	return blueskyPostCmd
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
