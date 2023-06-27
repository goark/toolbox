package webpage

import (
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"go.uber.org/zap"
)

// CreatePool creates pool and starts goroutine,
func (wp *Webpage) CreatePool() {
	if wp == nil {
		return
	}
	if wp.itemPool == nil {
		wp.itemPool = newItemPool()
		wp.itemPool.pool.start()
	}
}

// StopPool creates pool and starts goroutine,
func (wp *Webpage) StopPool() {
	if wp == nil || wp.itemPool == nil {
		return
	}
	wp.itemPool.done()
}

func (wp *Webpage) GetErrorInPool() error {
	if wp == nil {
		return nil
	}
	return wp.itemPool.errList.ErrorOrNil()
}

func (wp *Webpage) GetInfoInPool() []*Info {
	if wp == nil {
		return []*Info{}
	}
	list := wp.itemPool.getInfo()
	wp.Logger().Debug("GetInfoInPool", zap.Any("info", list))
	return list
}

func (wp *Webpage) SaveInfoToCache() []*Info {
	if wp == nil {
		return nil
	}
	wp.Logger().Info("save info data to cache")
	list := wp.GetInfoInPool()
	SortInfo(list)
	wp.cacheData.Puts(list...)
	return list
}

func (wp *Webpage) SaveCache() error {
	if wp == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	wp.Logger().Info("save webpage cache to file")
	return wp.cacheData.Save()
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
