package webpage

import (
	"sync"
)

type infoPool struct {
	mu   sync.RWMutex
	ch   chan *Webpage
	done chan struct{}
	pool []*Webpage
}

func newPool() *infoPool {
	return &infoPool{ch: make(chan *Webpage), done: make(chan struct{}, 1), pool: []*Webpage{}}
}

func (p *infoPool) put(page *Webpage) {
	if p == nil {
		return
	}
	p.ch <- page
}

func (p *infoPool) start() {
	if p == nil {
		return
	}
	go func() {
		for {
			page, ok := <-p.ch
			if !ok {
				break
			}
			func() {
				p.mu.Lock()
				defer p.mu.Unlock()
				p.pool = append(p.pool, page)
			}()
		}
		p.done <- struct{}{}
	}()
}

func (p *infoPool) length() int {
	if p == nil {
		return 0
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.pool)
}

func (p *infoPool) getList() []*Webpage {
	if p == nil {
		return []*Webpage{}
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	cpy := make([]*Webpage, len(p.pool), cap(p.pool))
	copy(cpy, p.pool)
	return cpy
}

func (p *infoPool) stop() {
	if p == nil {
		return
	}
	// time.Sleep(1 * time.Second)
	close(p.ch)
	<-p.done
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
