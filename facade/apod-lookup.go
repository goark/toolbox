package facade

import (
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/values"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newAPODLookupCmd returns cobra.Command instance for show sub-command
func newAPODLookupCmd(ui *rwi.RWI) *cobra.Command {
	apodLookupCmd := &cobra.Command{
		Use:     "lookup",
		Aliases: []string{"look", "l"},
		Short:   "Lookup APOD data NASA API key",
		Long:    "Lookup Astronomy Picture of the Day data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			apd, err := getAPOD()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			dateStr, err := cmd.Flags().GetString("date")
			if err != nil {
				return debugPrint(ui, err)
			}
			date, err := values.DateFrom(dateStr)
			if err != nil {
				return debugPrint(ui, err)
			}

			// lookup APOD data
			res, err := apd.Lookup(cmd.Context(), date)
			if err != nil {
				apd.Logger().Error("error in apod.Lookup", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			return debugPrint(ui, res.Encode(ui.Writer()))
		},
	}
	apodLookupCmd.Flags().StringP("date", "d", "", "Date for APOD data (YYYY-MM-DD)")

	return apodLookupCmd
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
