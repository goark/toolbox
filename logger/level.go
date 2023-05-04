package logger

import "github.com/rs/zerolog"

var levelMap = map[string]zerolog.Level{
	"nop":   zerolog.NoLevel,
	"error": zerolog.ErrorLevel,
	"warn":  zerolog.WarnLevel,
	"info":  zerolog.InfoLevel,
	"debug": zerolog.DebugLevel,
	"trace": zerolog.TraceLevel,
}

func LevelList() []string {
	return []string{"nop", "error", "warn", "info", "debug", "trace"}
}

func LevelFrom(s string) zerolog.Level {
	if lvl, ok := levelMap[s]; ok {
		return lvl
	}
	return zerolog.NoLevel
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
