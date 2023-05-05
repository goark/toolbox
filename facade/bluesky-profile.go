package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/spf13/cobra"
)

// newBlueskyProfileCmd returns cobra.Command instance for show sub-command
func newBlueskyProfileCmd(ui *rwi.RWI) *cobra.Command {
	blueskyProfileCmd := &cobra.Command{
		Use:     "profile",
		Aliases: []string{"prof"},
		Short:   "Output profile",
		Long:    "Show profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			bluesky, err := getBluesky()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			handle, err := cmd.Flags().GetString("handle")
			if err != nil {
				return debugPrint(ui, err)
			}
			jsonFlag, err := cmd.Flags().GetBool("json")
			if err != nil {
				return debugPrint(ui, err)
			}

			// post message
			if err := bluesky.ShowProfile(cmd.Context(), handle, jsonFlag, ui.Writer()); err != nil {
				bluesky.Logger().Error().Interface("error", errs.Wrap(err)).Send()
				return debugPrint(ui, err)
			}
			return nil
		},
	}
	blueskyProfileCmd.Flags().StringP("handle", "", "", "Handle name")
	blueskyProfileCmd.Flags().BoolP("json", "j", false, "Output JSON format")

	return blueskyProfileCmd
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
