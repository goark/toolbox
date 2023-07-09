package webpage

import (
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"go.uber.org/zap"
)

// CreatePool creates pool and starts goroutine,
func (cfg *Config) CreatePool() {
	if cfg == nil {
		return
	}
	if cfg.itemPool == nil {
		cfg.itemPool = newItemPool()
	}
}

// StopPool creates pool and starts goroutine,
func (cfg *Config) StopPool() ([]*Webpage, error) {
	if cfg == nil || cfg.itemPool == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	cfg.itemPool.done()
	if err := cfg.GetErrorInPool(); err != nil {
		return nil, errs.Wrap(err)
	}
	// get pages from pool
	return cfg.GetPagesFromPool(), nil
}

func (cfg *Config) GetErrorInPool() error {
	if cfg == nil {
		return nil
	}
	return cfg.itemPool.errList.ErrorOrNil()
}

func (cfg *Config) GetPagesFromPool() []*Webpage {
	if cfg == nil {
		return []*Webpage{}
	}
	list := cfg.itemPool.getPages()
	if len(list) > 1 {
		SortPages(list)
	}
	cfg.Logger().Debug("GetPagesFromPool", zap.Any("pages", list))
	return list
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
