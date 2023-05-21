package logger

import (
	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
	"github.com/ipfs/go-log/v2"
)

type Level int

const (
	LevelUnknown Level = iota
	LevelNop
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var levelMap = map[string]Level{
	"nop":   LevelNop,
	"error": LevelError,
	"warn":  LevelWarn,
	"info":  LevelInfo,
	"debug": LevelDebug,
	"trace": LevelTrace,
}

func LevelList() []string {
	return []string{"nop", "error", "warn", "info", "debug", "trace"}
}

func LevelFrom(s string) Level {
	if lvl, ok := levelMap[s]; ok {
		return lvl
	}
	return LevelUnknown
}

var gologLevelMap = map[Level]log.LogLevel{
	LevelError: log.LevelError,
	LevelWarn:  log.LevelWarn,
	LevelInfo:  log.LevelInfo,
	LevelDebug: log.LevelDebug,
	LevelTrace: log.LevelDebug,
}

func (lvl Level) GologLevel() (log.LogLevel, error) {
	if gologLvl, ok := gologLevelMap[lvl]; ok {
		return gologLvl, nil
	}
	return 0, errs.Wrap(ecode.ErrLogLevel)
}

func (lvl Level) String() string {
	for k, v := range levelMap {
		if lvl == v {
			return k
		}
	}
	return "unknown"
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
