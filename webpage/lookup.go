package webpage

import (
	"context"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
)

func (wp *Bookmark) Lookup(ctx context.Context, urlStr string, saveFlag bool) (*Info, error) {
	if wp == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	if info := wp.cacheData.Get(urlStr); info != nil {
		return info, nil
	}

	info, err := ReadPage(ctx, urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	wp.cacheData.Put(info)

	if saveFlag {
		if err := wp.cacheData.Save(); err != nil {
			return nil, errs.Wrap(err, errs.WithContext("url", urlStr), errs.WithContext("save", saveFlag))
		}
	}
	return info, nil
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
