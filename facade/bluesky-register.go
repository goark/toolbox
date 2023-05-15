package facade

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/bluesky"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/spf13/cobra"
)

// newBlueskyRegisterCmd returns cobra.Command instance for show sub-command
func newBlueskyRegisterCmd(ui *rwi.RWI) *cobra.Command {
	blueskyRegisterCmd := &cobra.Command{
		Use:     "register",
		Aliases: []string{"reg"},
		Short:   "Register account in local PC",
		Long:    "Register Bluesky account in local PC.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options (interactive mode)
			server, err := getBlueskyServer(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			handle, err := getBlueskyHandle(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			passaord, err := getBlueskyPassword(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			bcfg, err := bluesky.Register(cmd.Context(), server, handle, passaord, gopts.CacheDir, gopts.Logger)
			if err != nil {
				gopts.Logger.Error().Interface("error", errs.Wrap(err)).Send()
				return debugPrint(ui, err)
			}
			if err := debugPrint(ui, bcfg.Export(gopts.bskyConfigPath)); err != nil {
				return debugPrint(ui, err)
			}
			_ = ui.Outputln()
			_ = ui.Outputln("Host:", bcfg.Host)
			_ = ui.Outputln(" DID:", bcfg.Handle)
			_ = ui.Outputln()
			_ = ui.Outputln("output:", gopts.bskyConfigPath)
			return nil
		},
	}

	return blueskyRegisterCmd
}

func getBlueskyServer(ctx context.Context) (string, error) {
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return fmt.Fprintf(w, "     Host (default %s) > ", bluesky.DefaltHostName)
		},
	}
	text, err := editor.ReadLine(ctx)
	if err != nil {
		return "", errs.Wrap(err)
	}
	return strings.TrimSpace(text), nil
}

func getBlueskyHandle(ctx context.Context) (string, error) {
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return fmt.Fprint(w, "User (Handle/DID/email address) > ")
		},
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

func getBlueskyPassword(ctx context.Context) (string, error) {
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return fmt.Fprint(w, "                   App password > ")
		},
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
