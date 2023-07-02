package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/webpage"
	"github.com/spf13/cobra"
)

// newBookmarkCmd returns cobra.Command instance for show sub-command
func newWebpageCmd(ui *rwi.RWI) *cobra.Command {
	webpageCmd := &cobra.Command{
		Use:     "webpage",
		Aliases: []string{"web", "w", "bookmark", "book", "bm"},
		Short:   "Handling information for Web pages",
		Long:    "Handling information for Web pages.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	webpageCmd.PersistentFlags().StringP("url", "u", "", "Web page URL")
	_ = webpageCmd.MarkFlagRequired("url")
	webpageCmd.PersistentFlags().BoolP("save", "", false, "Save APOD data to cache")

	webpageCmd.AddCommand(
		newBookmarkLookupCmd(ui),
		newBookmarkPostCmd(ui),
	)
	return webpageCmd
}

func (gopts *globalOptions) getWebpage() (*webpage.Webpage, error) {
	cfg, err := webpage.New(gopts.CacheDir, gopts.Logger)
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
