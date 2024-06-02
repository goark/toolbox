package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/calendar"
	"github.com/goark/toolbox/ecode"
	"github.com/spf13/cobra"
)

// newCalendarCmd returns cobra.Command instance for show sub-command
func newCalendarCmd(ui *rwi.RWI) *cobra.Command {
	calendarCmd := &cobra.Command{
		Use:     "calendar",
		Aliases: []string{"cal", "c"},
		Short:   "Astronomical calendar commands",
		Long:    "Commands for Astronomical calendar by NAOJ https://eco.mtk.nao.ac.jp/koyomi/cande/calendar.html.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	calendarCmd.PersistentFlags().StringP("start", "", "", "start of date (YYYY-MM-DD)")
	calendarCmd.PersistentFlags().StringP("end", "", "", "end of date (YYYY-MM-DD)")
	calendarCmd.PersistentFlags().BoolP("holiday", "", false, "output holiday")
	calendarCmd.PersistentFlags().BoolP("ephemeris-all", "", false, "output all ephemeris")
	calendarCmd.PersistentFlags().BoolP("moon-phase", "", false, "output moon-phase")
	calendarCmd.PersistentFlags().BoolP("solar-term", "", false, "output solar-term")
	calendarCmd.PersistentFlags().BoolP("eclipse", "", false, "output eclipse")
	calendarCmd.PersistentFlags().BoolP("planet", "", false, "output planet")
	calendarCmd.PersistentFlags().StringP("template", "", "", "template file for Output format")

	calendarCmd.AddCommand(
		newCalendarLookupCmd(ui),
		newCalendarPostCmd(ui),
	)
	return calendarCmd
}

func getCalendarConfig(cmd *cobra.Command, gopts *globalOptions) (*calendar.Config, error) {
	start, err := cmd.Flags().GetString("start")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	end, err := cmd.Flags().GetString("end")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	holidayFlag, err := cmd.Flags().GetBool("holiday")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	moonPhaseFlag, err := cmd.Flags().GetBool("moon-phase")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	solarTermFlag, err := cmd.Flags().GetBool("solar-term")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	eclipseFlag, err := cmd.Flags().GetBool("eclipse")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	planetFlag, err := cmd.Flags().GetBool("planet")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	ephemerisAllFlag, err := cmd.Flags().GetBool("ephemeris-all")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if ephemerisAllFlag {
		moonPhaseFlag = true
		solarTermFlag = true
		eclipseFlag = true
		planetFlag = true
	}
	templateFile, err := cmd.Flags().GetString("template")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return calendar.NewConfig(start, end, holidayFlag, moonPhaseFlag, solarTermFlag, eclipseFlag, planetFlag, gopts.TempDir, templateFile)
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
