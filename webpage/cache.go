package webpage

import (
	"sync"
)

// Cache is cache data for Web page informations.
type Cache struct {
	mu sync.RWMutex
	// dir   string
	pages map[string]*Webpage
}

// NewCache function creates Cache instance and import cache data.
func NewCache(dir string) *Cache {
	return &Cache{pages: map[string]*Webpage{}}
}

// Get method gets Web page data from cache.
func (c *Cache) Get(urlStr string) *Webpage {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pages[urlStr]
}

// Put method puts Web page data to cache.
func (c *Cache) Put(page *Webpage) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pages == nil {
		c.pages = map[string]*Webpage{}
	}
	if page != nil && len(page.URL) > 0 {
		c.pages[page.URL] = page
	}
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
