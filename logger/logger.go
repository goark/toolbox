package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goark/errs"
	"github.com/goark/gocli/cache"
	"github.com/rs/zerolog"
)

// New function returns new zerolog.Logger instance.
func New(lvl zerolog.Level, dir string) (*zerolog.Logger, error) {
	logger := zerolog.Nop()
	if lvl == zerolog.NoLevel {
		return &logger, nil
	}
	logpath := getPath(dir)
	if file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		return &logger, errs.Wrap(err, errs.WithContext("logpath", logpath))
	} else {
		logger = zerolog.New(file).Level(lvl).With().Timestamp().Logger()
	}
	return &logger, nil
}

// DefaultLogDir function returns default log directory ($XDG_CACHE_HOME/appName/)
func DefaultLogDir(appName string) string {
	return cache.Dir(appName)
}

func getPath(dir string) string {
	if len(dir) == 0 {
		dir = "."
	}
	_ = os.MkdirAll(dir, 0700)
	return filepath.Join(dir, fmt.Sprintf("access.%s.log", time.Now().Local().Format("20060102")))
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
