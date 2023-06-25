package facade

import (
	"encoding/json"

	"github.com/goark/gocli/rwi"
	"github.com/spf13/cobra"
)

// newFeedLookupCmd returns cobra.Command instance for show sub-command
func newFeedLookupCmd(ui *rwi.RWI) *cobra.Command {
	feedLookupCmd := &cobra.Command{
		Use:     "lookup",
		Aliases: []string{"look", "l"},
		Short:   "Lookup information for Web page",
		Long:    "Lookup information for Web page.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			cfg, err := gopts.getBookmark()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			saveFlag, err := cmd.Flags().GetBool("save")
			if err != nil {
				return debugPrint(ui, err)
			}

			// lookup feed
			list, err := getFeedAll(cmd, cfg)
			if err != nil {
				return debugPrint(ui, err)
			}
			if saveFlag {
				if err := cfg.SaveCache(); err != nil {
					return debugPrint(ui, err)
				}
			}
			return debugPrint(ui, json.NewEncoder(ui.Writer()).Encode(list))
		},
	}
	return feedLookupCmd
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
