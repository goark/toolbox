package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/cache"
	"github.com/goark/toolbox/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type globalOptions struct {
	Logger   *zerolog.Logger
	CacheDir string
}

func getGlobalOptions() (*globalOptions, error) {
	cacheDir := viper.GetString("cache-dir")
	if len(cacheDir) == 0 {
		cacheDir = cache.Dir(Name)
	}
	logger, err := logger.New(
		logger.LevelFrom(viper.GetString("log-level")),
		viper.GetString("log-dir"),
	)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &globalOptions{
		Logger:   logger,
		CacheDir: cacheDir,
	}, nil
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
