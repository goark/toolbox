package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/cache"
	"github.com/goark/gocli/config"
	"github.com/goark/toolbox/logger"
	"github.com/goark/toolbox/tempdir"
	"github.com/ipfs/go-log/v2"
	"github.com/spf13/viper"
)

type globalOptions struct {
	Logger          *log.ZapEventLogger
	CacheDir        string
	TempDir         *tempdir.TempDir
	bskyConfigPath  string
	mstdnConfigPath string
	apodConfigPath  string
}

func getGlobalOptions() (*globalOptions, error) {
	cacheDir := viper.GetString("cache-dir")
	if len(cacheDir) == 0 {
		cacheDir = cache.Dir(Name)
	}
	tempDir := tempdir.New(viper.GetString("temp-dir"))
	golog, err := logger.New(
		logger.LevelFrom(viper.GetString("log-level")),
		Name,
		config.Dir(Name),
		viper.GetString("log-dir"),
	)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	bskyConfigPath := viper.GetString("bluesky-config")
	if len(bskyConfigFile) == 0 {
		bskyConfigPath = defaultBskyConfigPath
	}
	mstdnConfigPath := viper.GetString("mastodon-config")
	if len(mstdnConfigPath) == 0 {
		mstdnConfigPath = defaultMstdnConfigPath
	}
	apodConfigPath := viper.GetString("apod-config")
	if len(mstdnConfigPath) == 0 {
		mstdnConfigPath = defaultAPODConfigPath
	}
	return &globalOptions{
		Logger:          golog,
		CacheDir:        cacheDir,
		TempDir:         tempDir,
		bskyConfigPath:  bskyConfigPath,
		mstdnConfigPath: mstdnConfigPath,
		apodConfigPath:  apodConfigPath,
	}, nil
}

/* Copyright 2023-2024 Spiegel
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
