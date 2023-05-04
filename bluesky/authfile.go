package bluesky

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluesky-social/indigo/xrpc"
	"github.com/goark/errs"
)

func (cfg *Bluesky) authPath() string {
	ss := strings.Split(cfg.Handle, ":")
	if len(ss) > 0 && len(ss[len(ss)-1]) > 0 {
		return filepath.Join(cfg.BaseDir(), ss[len(ss)-1]+".auth")
	}
	return filepath.Join(cfg.BaseDir(), "bluesky.auth")
}

func (cfg *Bluesky) readAuth() (*xrpc.AuthInfo, error) {
	file, err := os.Open(cfg.authPath())
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("authfile", cfg.authPath()))
	}
	defer file.Close()

	var auth xrpc.AuthInfo
	if err := json.NewDecoder(file).Decode(&auth); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("authfile", cfg.authPath()))
	}
	return &auth, nil
}

func (cfg *Bluesky) writeAuth(auth *xrpc.AuthInfo) error {
	file, err := os.Create(cfg.authPath())
	if err != nil {
		return errs.Wrap(err, errs.WithContext("authfile", cfg.authPath()))
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(auth); err != nil {
		return errs.Wrap(err, errs.WithContext("authfile", cfg.authPath()))
	}
	return nil
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
