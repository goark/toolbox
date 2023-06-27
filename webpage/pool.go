package webpage

import "sync"

type infoPool struct {
	mu   sync.RWMutex
	ch   chan *Info
	done chan struct{}
	pool []*Info
}

func newPool() *infoPool {
	return &infoPool{ch: make(chan *Info), done: make(chan struct{}, 1), pool: []*Info{}}
}

func (p *infoPool) put(i *Info) {
	if p == nil {
		return
	}
	p.ch <- i
}

func (p *infoPool) start() {
	if p == nil {
		return
	}
	go func() {
		for {
			i, ok := <-p.ch
			if !ok {
				break
			}
			func() {
				p.mu.Lock()
				defer p.mu.Unlock()
				p.pool = append(p.pool, i)
			}()
		}
		p.done <- struct{}{}
	}()
}

func (p *infoPool) getList() []*Info {
	if p == nil {
		return []*Info{}
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	cpy := make([]*Info, len(p.pool), cap(p.pool))
	copy(cpy, p.pool)
	return cpy
}

func (p *infoPool) stop() {
	if p == nil {
		return
	}
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
