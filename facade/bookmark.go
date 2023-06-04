package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/bookmark"
	"github.com/goark/toolbox/ecode"
	"github.com/spf13/cobra"
)

// newBookmarkCmd returns cobra.Command instance for show sub-command
func newBookmarkCmd(ui *rwi.RWI) *cobra.Command {
	bookmarkCmd := &cobra.Command{
		Use:     "bookmark",
		Aliases: []string{"book", "bm"},
		Short:   "Handling information for Web pages",
		Long:    "Handling information for Web pages.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	bookmarkCmd.PersistentFlags().StringP("url", "u", "", "Web page URL")
	_ = bookmarkCmd.MarkFlagRequired("url")
	bookmarkCmd.PersistentFlags().BoolP("save", "", false, "Save APOD data to cache")

	bookmarkCmd.AddCommand(
		newBookmarkDLookupCmd(ui),
		newBookmarkPostCmd(ui),
	)
	return bookmarkCmd
}

func (gopts *globalOptions) getBookmark() (*bookmark.Config, error) {
	cfg, err := bookmark.New(gopts.CacheDir, gopts.Logger)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return cfg, nil
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
