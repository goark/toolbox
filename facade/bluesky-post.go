package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/spf13/cobra"
)

// newVersionCmd returns cobra.Command instance for show sub-command
func newBlueskyPostCmd(ui *rwi.RWI) *cobra.Command {
	blueskyPostCmd := &cobra.Command{
		Use:     "post",
		Aliases: []string{"pst", "p"},
		Short:   "Post message to Bluesky",
		Long:    "Post message to Bluesky.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			bluesky, err := getBluesky()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			msg, err := cmd.Flags().GetString("message")
			if err != nil {
				return debugPrint(ui, err)
			}

			// post message
			resText, err := bluesky.PostMessage(cmd.Context(), msg)
			if err != nil {
				bluesky.Logger().Error().Interface("error", errs.Wrap(err)).Send()
				return debugPrint(ui, err)
			}
			return debugPrint(ui, ui.Outputln(resText))
		},
	}
	blueskyPostCmd.Flags().StringP("message", "m", "", "Message")

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
