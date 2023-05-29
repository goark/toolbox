package facade

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/apod"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/spf13/cobra"
)

// newAPODRegisterCmd returns cobra.Command instance for show sub-command
func newAPODRegisterCmd(ui *rwi.RWI) *cobra.Command {
	apodRegisterCmd := &cobra.Command{
		Use:     "register",
		Aliases: []string{"reg"},
		Short:   "Register NASA API key",
		Long:    "Register NASA API key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options (interactive mode)
			apiKey, err := getAPODAPIKey(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			mcfg := apod.Register(apiKey, gopts.CacheDir, gopts.Logger)
			if err := debugPrint(ui, mcfg.Export(gopts.apodConfigPath)); err != nil {
				return debugPrint(ui, err)
			}
			return debugPrint(ui, ui.Outputln("output:", gopts.apodConfigPath))
		},
	}

	return apodRegisterCmd
}

func getAPODAPIKey(ctx context.Context) (string, error) {
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) { return fmt.Fprint(w, "NASA API key > ") },
	}
	for {
		text, err := editor.ReadLine(ctx)
		if err != nil {
			return "", errs.Wrap(err)
		}
		text = strings.TrimSpace(text)
		if len(text) > 0 {
			return text, nil
		}
	}
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
