package webpage

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/logger"
	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
)

// Webpage is configuration for bookmark
type Webpage struct {
	cacheDir  string
	cacheData *Cache
	itemPool  *itemPool
	logger    *log.ZapEventLogger
}

// New functions creates new Config instance.
func New(cacheDir string, logger *log.ZapEventLogger) (*Webpage, error) {
	data, err := NewCache(cacheDir)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("cache_dir", cacheDir))
	}
	return &Webpage{
		cacheDir:  cacheDir,
		cacheData: data,
		logger:    logger,
	}, nil
}

// Logger method returns zap.Logger instance.
func (wp *Webpage) Logger() *zap.Logger {
	if wp == nil || wp.logger == nil {
		return logger.Nop().Desugar()
	}
	return wp.logger.Desugar()
}

// PutURLToCache method puts Info of web page, and returns Info.
func (wp *Webpage) PutURLToCache(ctx context.Context, urlStr string) (*Info, error) {
	if wp == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	return wp.cacheData.PutURL(ctx, urlStr)
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
