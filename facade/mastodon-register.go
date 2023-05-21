package facade

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/mastodon"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// newBlueskyCmd returns cobra.Command instance for show sub-command
func newMastodonRegisterCmd(ui *rwi.RWI) *cobra.Command {
	mastodonRegisterCmd := &cobra.Command{
		Use:     "register",
		Aliases: []string{"reg"},
		Short:   "Register application",
		Long:    "Register Mastodon application.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// global options
			gopts, err := getGlobalOptions()
			if err != nil {
				return debugPrint(ui, err)
			}
			// local options (interactive mode)
			server, err := getMastodonServer(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			username, err := getMastodonUserId(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			passaord, err := getMastodonPassword(cmd.Context())
			if err != nil {
				return debugPrint(ui, err)
			}
			mcfg, err := mastodon.Register(cmd.Context(), server, username, passaord, gopts.Logger)
			if err != nil {
				gopts.Logger.Desugar().Error("error in mastodon.Register", zap.Object("error", zapobject.New(err)))
				return debugPrint(ui, err)
			}
			if err := debugPrint(ui, mcfg.Export(gopts.mstdnConfigPath)); err != nil {
				return debugPrint(ui, err)
			}
			_ = ui.Outputln()
			_ = ui.Outputln("          server:", mcfg.Server)
			_ = ui.Outputln("application name:", mcfg.AppName())
			_ = ui.Outputln("         website:", mcfg.Registory())
			_ = ui.Outputln("          scopes:", mcfg.Scopes())
			_ = ui.Outputln()
			_ = ui.Outputln("output:", gopts.mstdnConfigPath)
			return nil
		},
	}

	return mastodonRegisterCmd
}

func getMastodonServer(ctx context.Context) (string, error) {
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) { return fmt.Fprint(w, "Server (e.g. mastodon.social) > ") },
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

func getMastodonUserId(ctx context.Context) (string, error) {
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) { return fmt.Fprint(w, "         User (email address) > ") },
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

func getMastodonPassword(ctx context.Context) (string, error) {
	editor := readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) { return fmt.Fprint(w, "                     Password > ") },
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
