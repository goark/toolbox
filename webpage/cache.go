package webpage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
)

const (
	cachesFile = "webpage.cache.json"
)

// Cache is cache data for Web page informations.
type Cache struct {
	mu   sync.RWMutex
	dir  string
	info map[string]*Info
}

// NewCache function creates Cache instance and import cache data.
func NewCache(dir string) (*Cache, error) {
	cache := &Cache{dir: dir, info: map[string]*Info{}}
	file, err := os.Open(cache.path())
	if err != nil {
		return cache, nil
	}
	defer file.Close()

	var infos []*Info
	if err := json.NewDecoder(file).Decode(&infos); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("path", cache.path()))
	}

	for _, i := range infos {
		cache.info[i.URL] = i
	}
	return cache, nil
}

// Save method saves data to cache file.
func (c *Cache) Save() error {
	if c == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	if len(c.info) == 0 {
		c.Remove()
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	var infos []*Info
	for _, v := range c.info {
		infos = append(infos, v)
	}

	file, err := os.OpenFile(c.path(), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return errs.Wrap(err, errs.WithContext("path", c.path()))
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(infos); err != nil {
		return errs.Wrap(err, errs.WithContext("path", c.path()))
	}
	return nil
}

// Remove method removes cache file.
func (c *Cache) Remove() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = os.Remove(c.path())
}

// Get method gets Web page data from cache.
func (c *Cache) Get(urlStr string) *Info {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.info[urlStr]
}

// Put method puts Web page data to cache.
func (c *Cache) Put(i *Info) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.info == nil {
		c.info = map[string]*Info{}
	}
	if i != nil && len(i.URL) > 0 {
		c.info[i.URL] = i
	}
}

// Puts method puts list of Web page data to cache.
func (c *Cache) Puts(is ...*Info) {
	for _, i := range is {
		c.Put(i)
	}
}

// PutURL reads Web page data from URL and puts to cache.
func (c *Cache) PutURL(ctx context.Context, urlStr string) (*Info, error) {
	if c == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	if info := c.Get(urlStr); info != nil {
		return info, nil
	}
	info, err := ReadPage(ctx, urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	c.Put(info)
	return info, nil
}

func (c *Cache) path() string {
	if c == nil {
		return ""
	}
	return filepath.Join(c.dir, cachesFile)
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
