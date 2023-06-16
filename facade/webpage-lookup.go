package facade

import (
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newBookmarkDLookupCmd returns cobra.Command instance for show sub-command
func newBookmarkDLookupCmd(ui *rwi.RWI) *cobra.Command {
	bookmarkLookupCmd := &cobra.Command{
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
			urlStr, err := cmd.Flags().GetString("url")
			if err != nil {
				return debugPrint(ui, err)
			}
			saveFlag, err := cmd.Flags().GetBool("save")
			if err != nil {
				return debugPrint(ui, err)
			}

			// lookup Web page data
			info, err := cfg.Lookup(cmd.Context(), urlStr, saveFlag)
			if err != nil {
				gopts.Logger.Desugar().Error("error in bookmark.Lookup", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			return debugPrint(ui, info.Encode(ui.Writer()))
		},
	}
	return bookmarkLookupCmd
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
