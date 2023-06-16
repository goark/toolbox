package webpage

import (
	"fmt"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/toolbox/logger"
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
)

// Config is configuration for bookmark
type Bookmark struct {
	cacheDir  string
	cacheData *Cache
	logger    *log.ZapEventLogger
}

// New functions creates new Config instance.
func New(cacheDir string, logger *log.ZapEventLogger) (*Bookmark, error) {
	data, err := NewCache(cacheDir)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("cache_dir", cacheDir))
	}
	return &Bookmark{
		cacheDir:  cacheDir,
		cacheData: data,
		logger:    logger,
	}, nil
}

func (info *Info) MakeMessage(prefixMsg string) string {
	if info == nil {
		return ""
	}
	bld := strings.Builder{}

	//title
	if len(info.Title) > 0 {
		bld.WriteString(fmt.Sprintln(prefixMsg, info.Title))
	}
	// URL
	bld.WriteString(fmt.Sprintln(info.URL))
	return bld.String()
}

// Logger method returns zap.Logger instance.
func (wp *Bookmark) Logger() *zap.Logger {
	if wp == nil || wp.logger == nil {
		return logger.Nop().Desugar()
	}
	return wp.logger.Desugar()
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
