package facade

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/apod"
	"github.com/goark/toolbox/ecode"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newAPODCmd returns cobra.Command instance for show sub-command
func newAPODCmd(ui *rwi.RWI) *cobra.Command {
	apodCmd := &cobra.Command{
		Use:     "apod",
		Aliases: []string{},
		Short:   "NASA APOD commands",
		Long:    "Commands for Astronomy Picture of the Day by NASA API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	apodCmd.PersistentFlags().StringP("date", "d", "", "Date for APOD data (YYYY-MM-DD)")
	apodCmd.PersistentFlags().BoolP("utc", "u", false, "Time base on UTC")

	apodCmd.AddCommand(
		newAPODRegisterCmd(ui),
		newAPODLookupCmd(ui),
		newAPODPostCmd(ui),
	)
	return apodCmd
}

func (gopts *globalOptions) getAPOD(ctx context.Context) (*apod.APOD, error) {
	acfg, err := apod.New(ctx, gopts.apodConfigPath, gopts.CacheDir, gopts.Logger)
	if err != nil {
		err = errs.Wrap(err)
		gopts.Logger.Desugar().Error("cannot get configuration for Mastodon", zap.Object("error", zapobject.New(err)))
		return nil, err
	}
	return acfg, nil
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
