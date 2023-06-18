package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/ecode"
	"github.com/spf13/cobra"
)

// newFeedCmd returns cobra.Command instance for show sub-command
func newFeedCmd(ui *rwi.RWI) *cobra.Command {
	webpageCmd := &cobra.Command{
		Use:     "feed",
		Aliases: []string{"rss"},
		Short:   "Handling information for Web feed",
		Long:    "Handling information for Web feed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	webpageCmd.PersistentFlags().StringP("url", "u", "", "Feed URL")
	webpageCmd.PersistentFlags().StringP("flickr-id", "", "", "Flickr ID")
	webpageCmd.MarkFlagsMutuallyExclusive("url", "flickr-id")
	webpageCmd.PersistentFlags().BoolP("save", "", false, "Save webpage data to cache")

	webpageCmd.AddCommand(
		newFeedLookupCmd(ui),
		newFeedPostCmd(ui),
	)
	return webpageCmd
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
