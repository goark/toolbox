package webpage

import (
	"errors"
	"sync"
)

type errorList struct {
	mu      sync.Mutex
	errList []error
}

func newErrorList() *errorList {
	return &errorList{errList: []error{}}
}

func (el *errorList) Add(err error) {
	if el == nil {
		return
	}
	el.mu.Lock()
	defer el.mu.Unlock()
	el.errList = append(el.errList, err)
}

func (el *errorList) GetError() error {
	if el == nil {
		return nil
	}
	el.mu.Lock()
	defer el.mu.Unlock()
	if len(el.errList) == 0 {
		return nil
	}
	return errors.Join(el.errList...)
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