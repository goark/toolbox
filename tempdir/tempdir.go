package tempdir

import (
	"os"
	"path/filepath"

	"github.com/goark/errs"
	"github.com/goark/toolbox/consts"
	"github.com/goark/toolbox/ecode"
)

// TempDir type is temporary directory class.
type TempDir struct {
	baseDir   string
	pattern   string
	makedPath string
}

// New function creates new TempDir instance.
func New(baseDir string) *TempDir {
	if len(baseDir) == 0 {
		baseDir = os.TempDir()
	}
	return &TempDir{baseDir: baseDir, pattern: consts.AppNameShort + "_*", makedPath: ""}
}

// String is stringer method.
func (t *TempDir) String() string {
	if t == nil {
		return ""
	}
	if len(t.makedPath) == 0 {
		return t.baseDir
	}
	return t.makedPath
}

// Path method returns path of temporary directory.
func (t *TempDir) Path() string {
	if t == nil || len(t.makedPath) == 0 {
		return ""
	}
	return t.makedPath
}

// MakeDir method makes temporary directory.
func (t *TempDir) MakeDir() error {
	if t == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	if len(t.makedPath) > 0 {
		return nil
	}
	dir, err := os.MkdirTemp(t.baseDir, t.pattern)
	if err == nil {
		t.makedPath = dir
	}
	return err
}

// MakeDir method removes temporary directory.
func (t *TempDir) CleanUp() error {
	if t == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	if len(t.makedPath) == 0 {
		return nil
	}
	return os.RemoveAll(t.makedPath)
}

// FilePath method makes path of file in temporary directory.
func (t *TempDir) FilePath(filename string) string {
	if t == nil {
		return ""
	}
	return filepath.Join(t.Path(), filename)
}

/* Copyright 2024 Spiegel
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
