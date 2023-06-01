package facade

import (
	"github.com/goark/gocli/rwi"
	"github.com/spf13/cobra"
)

// newBlueskyCmd returns cobra.Command instance for show sub-command
func newMastodonProfileCmd(ui *rwi.RWI) *cobra.Command {
	mastodonProfileCmd := &cobra.Command{
		Use:     "profile",
		Aliases: []string{"prof"},
		Short:   "Output my profile",
		Long:    "Output my profile.",
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
			jsonFlag, err := cmd.Flags().GetBool("json")
			if err != nil {
				return debugPrint(ui, err)
			}

			// get my account
			return debugPrint(ui, mstdn.ShowProfile(cmd.Context(), jsonFlag, ui.Writer()))
		},
	}
	mastodonProfileCmd.Flags().BoolP("json", "j", false, "Output JSON format")

	return mastodonProfileCmd
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
