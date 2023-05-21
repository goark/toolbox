package facade

import (
	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/mastodon"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newBlueskyCmd returns cobra.Command instance for show sub-command
func newMastodonCmd(ui *rwi.RWI) *cobra.Command {
	mastodonCmd := &cobra.Command{
		Use:     "mastodon",
		Aliases: []string{"mstdn", "mast", "mst"},
		Short:   "Simple Mastodon commands",
		Long:    "Simple Mastodon commands.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	mastodonCmd.AddCommand(
		newMastodonRegisterCmd(ui),
		newMastodonProfileCmd(ui),
		newMastodonPostCmd(ui),
	)
	return mastodonCmd
}

func getMastodon() (*mastodon.Mastodon, error) {
	gopts, err := getGlobalOptions()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mcfg, err := mastodon.New(gopts.mstdnConfigPath, gopts.Logger)
	if err != nil {
		err = errs.Wrap(err)
		gopts.Logger.Desugar().Error("cannot get configuration for Mastodon", zap.Object("error", zapobject.New(err)))
		return nil, err
	}
	return mcfg, nil
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
