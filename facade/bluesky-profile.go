package facade

import (
	"io"
	"strings"

	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newBlueskyProfileCmd returns cobra.Command instance for show sub-command
func newBlueskyProfileCmd(ui *rwi.RWI) *cobra.Command {
	blueskyProfileCmd := &cobra.Command{
		Use:     "profile",
		Aliases: []string{"prof"},
		Short:   "Output Bluesky profile",
		Long:    "Output Bluesky profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			bsky, err := gopts.getBluesky()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options
			jsonFlag, err := cmd.Flags().GetBool("json")
			if err != nil {
				return debugPrint(ui, err)
			}
			handle, err := cmd.Flags().GetString("handle")
			if err != nil {
				return debugPrint(ui, err)
			}
			pipeFlag, err := cmd.Flags().GetBool("pipe")
			if err != nil {
				return debugPrint(ui, err)
			}
			if pipeFlag {
				b, err := io.ReadAll(ui.Reader())
				if err != nil {
					return debugPrint(ui, err)
				}
				handle = string(b)
			}
			handle = strings.TrimSpace(handle)

			// post message
			if err := bsky.ShowProfile(cmd.Context(), handle, jsonFlag, ui.Writer()); err != nil {
				bsky.Logger().Error("error in bluesky.ShowProfile", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			return nil
		},
	}
	blueskyProfileCmd.Flags().BoolP("json", "j", false, "Output JSON format")
	blueskyProfileCmd.Flags().StringP("handle", "", "", "Handle name")
	blueskyProfileCmd.Flags().BoolP("pipe", "", false, "Input from standard-input")
	blueskyProfileCmd.MarkFlagsMutuallyExclusive("handle", "pipe")

	return blueskyProfileCmd
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
