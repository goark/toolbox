package facade

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/gocli/cache"
	"github.com/goark/gocli/config"
	"github.com/goark/gocli/exitcode"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/consts"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//Name is applicatin name
	Name = consts.AppNameShort
	//Version is version for applicatin
	Version = "dev-version"
)

const (
	configFile        = "config"
	bskyConfigFile    = "bluesky.json"
	mstdnConfigFile   = "mastodon.json"
	nasaapiConfigFile = "nasaapi.json"
)

var (
	debugFlag              bool   //debug flag
	configPath             string //path for config file
	defaultConfigPath      = config.Path(Name, configFile+".yaml")
	defaultBskyConfigPath  = config.Path(Name, bskyConfigFile)
	defaultMstdnConfigPath = config.Path(Name, mstdnConfigFile)
	defaultAPODConfigPath  = config.Path(Name, nasaapiConfigFile)
)

// newRootCmd returns cobra.Command instance for root command
func newRootCmd(ui *rwi.RWI, args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   Name,
		Short: "A collection of miscellaneous commands",
		Long:  "A collection of miscellaneous commands.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	// global options (binding)
	rootCmd.PersistentFlags().StringP("cache-dir", "", cache.Dir(Name), "Directory for cache files")
	rootCmd.PersistentFlags().StringP("log-dir", "", logger.DefaultLogDir(Name), "Directory for log files")
	rootCmd.PersistentFlags().StringP("log-level", "", "nop", fmt.Sprintf("Log level [%s]", strings.Join(logger.LevelList(), "|")))
	rootCmd.PersistentFlags().StringP("bluesky-config", "", defaultBskyConfigPath, "Config file for Bluesky")
	rootCmd.PersistentFlags().StringP("mastodon-config", "", defaultMstdnConfigPath, "Config file for Mastodon")
	rootCmd.PersistentFlags().StringP("apod-config", "", defaultAPODConfigPath, "Config file for APOD")

	//Bind config file
	_ = viper.BindPFlag("cache-dir", rootCmd.PersistentFlags().Lookup("log-dir"))
	_ = viper.BindPFlag("log-dir", rootCmd.PersistentFlags().Lookup("log-dir"))
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("bluesky-config", rootCmd.PersistentFlags().Lookup("bluesky-config"))
	_ = viper.BindPFlag("mastodon-config", rootCmd.PersistentFlags().Lookup("mastodon-config"))
	_ = viper.BindPFlag("apod-config", rootCmd.PersistentFlags().Lookup("apod-config"))
	cobra.OnInitialize(initConfig)

	// global options (other)
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "", false, "for debug")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "", defaultConfigPath, "Config file")

	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetArgs(args)
	rootCmd.SetIn(ui.Reader())       //Stdin
	rootCmd.SetOut(ui.ErrorWriter()) //Stdout -> Stderr
	rootCmd.SetErr(ui.ErrorWriter()) //Stderr
	rootCmd.AddCommand(
		newVersionCmd(ui),
		newBlueskyCmd(ui),
		newMastodonCmd(ui),
		newAPODCmd(ui),
		newWebpageCmd(ui),
		newFeedCmd(ui),
	)
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configPath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configPath)
		// configDir = filepath.Dir(configPath)
	} else {
		// Find config directory.
		configDir := config.Dir(Name)
		if len(configDir) == 0 {
			configDir = "." //current directory
		}
		// Search config in home directory with name "config.yaml" (without extension).
		viper.AddConfigPath(configDir)
		viper.SetConfigName(configFile)
	}
	viper.AutomaticEnv()     // read in environment variables that match
	_ = viper.ReadInConfig() // If a config file is found, read it in.
}

// Execute is called from main function
func Execute(ui *rwi.RWI, args []string) (exit exitcode.ExitCode) {
	defer func() {
		//panic hundling
		if r := recover(); r != nil {
			_ = ui.OutputErrln("Panic:", r)
			for depth := 0; ; depth++ {
				pc, _, line, ok := runtime.Caller(depth)
				if !ok {
					break
				}
				_ = ui.OutputErrln(" ->", depth, ":", runtime.FuncForPC(pc).Name(), ": line", line)
			}
			exit = exitcode.Abnormal
		}
	}()

	// create interrupt SIGNAL
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//execution
	exit = exitcode.Normal
	if err := newRootCmd(ui, args).ExecuteContext(ctx); err != nil {
		exit = exitcode.Abnormal
	}
	return
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
