package facade

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/consts"
	"github.com/spf13/cobra"
)

// var versionStrings = []string{ //output message of version
// 	Name + " " + Version,
// 	"repository: " + consts.RepositoryURL,
// }

// newVersionCmd returns cobra.Command instance for show sub-command
func newVersionCmd(ui *rwi.RWI) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver", "v"},
		Short:   "Print the version number",
		Long:    "Print the version number of " + Name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ui.OutputErrln(getVersion())
		},
	}

	return versionCmd
}

// getVersion returns a formatted version message including name, version, and repository URL.
func getVersion() string {
	return strings.Join(
		[]string{
			strings.Join([]string{Name, replaceVersion(Version)}, " "),
			fmt.Sprintf("repository: %v", consts.RepositoryURL),
		},
		"\n",
	)
}

// replaceVersion returns the version string based on the provided version and build information
func replaceVersion(v string) string {
	if v != "" {
		return v
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion := fmt.Sprintf("(compiled with %v)", info.GoVersion)
		if info.Main.Version != "" {
			return joinNonEmpty(info.Main.Version, goVersion)
		}
		var revision, dirty string
		for _, v := range info.Settings {
			switch v.Key {
			case "vcs.revision":
				revision = v.Value
			case "vcs.modified":
				if v.Value == "true" {
					dirty = "(dirty)"
				}
			}
		}
		if revision != "" {
			return joinNonEmpty(revision, dirty, goVersion)
		}
	}
	return "(version not set)"
}

func joinNonEmpty(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, " ")
}

/* Copyright 2023-2026 Spiegel
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
