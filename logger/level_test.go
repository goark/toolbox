package logger

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestDate(t *testing.T) {
	want := strings.Join([]string{"nop", "error", "warn", "info", "debug", "trace"}, "|")
	if got := strings.Join(LevelList(), "|"); got != want {
		t.Errorf("LevelList() is \"%v\", want \"%v\"", got, want)
	}
}

func TestLevelFrom(t *testing.T) {
	testCases := []struct {
		s    string
		want zerolog.Level
	}{
		{s: "", want: zerolog.NoLevel},
		{s: "error", want: zerolog.ErrorLevel},
		{s: "warn", want: zerolog.WarnLevel},
		{s: "info", want: zerolog.InfoLevel},
		{s: "debug", want: zerolog.DebugLevel},
		{s: "trace", want: zerolog.TraceLevel},
		{s: "foo", want: zerolog.NoLevel},
	}
	for _, tc := range testCases {
		if got := LevelFrom(tc.s); got != tc.want {
			t.Errorf("LevelFrom(\"%s\") is [%v], want [%v]", tc.s, got, tc.want)
		}
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
