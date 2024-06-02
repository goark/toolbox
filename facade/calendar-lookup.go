package facade

import (
	"github.com/goark/gocli/rwi"
	"github.com/spf13/cobra"
)

// newCalendarLookupCmd returns cobra.Command instance for show sub-command
func newCalendarLookupCmd(ui *rwi.RWI) *cobra.Command {
	calendarLookupCmd := &cobra.Command{
		Use:     "lookup",
		Aliases: []string{"look", "l"},
		Short:   "Lookup astronomical calendar",
		Long:    "Lookup astronomical calendar.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			ccfg, err := getCalendarConfig(cmd, gopts)
			if err != nil {
				return debugPrint(ui, err)
			}

			// lookup calendar data
			return debugPrint(ui, ccfg.OutputEvent(ui.Writer()))
		},
	}

	return calendarLookupCmd
}

/* Copyright 2024 Spiegel
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
