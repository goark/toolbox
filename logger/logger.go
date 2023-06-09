package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goark/errs"
	"github.com/goark/gocli/cache"
	"github.com/ipfs/go-log/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// New function returns new zerolog.Logger instance.
func New(lvl Level, name, confDir, logDir string) (*log.ZapEventLogger, error) {
	_ = godotenv.Load(".env", filepath.Join(confDir, "env"))
	if lvl == LevelNop {
		return Nop(), nil
	}
	level, err := lvl.GologLevel()
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("level", lvl.String()))
	}
	cfg := log.GetConfig()
	cfg.Format = log.JSONOutput
	cfg.Stderr = false
	cfg.Stdout = false
	delete(cfg.Labels, name)
	cfg.Level = level
	cfg.File = getPath(logDir)
	log.SetupLogging(cfg)

	logger := log.Logger(name)
	logger.SugaredLogger = *(logger.Desugar().WithOptions(zap.WithCaller(false)).Sugar())
	return logger, nil
}

// Nop function returns nop logger
func Nop() *log.ZapEventLogger {
	logger := log.Logger("")
	logger.SugaredLogger = *zap.NewNop().Sugar()
	return logger
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
