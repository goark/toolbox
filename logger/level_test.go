package logger

import (
	"errors"
	"strings"
	"testing"

	"github.com/goark/toolbox/ecode"
	"github.com/ipfs/go-log/v2"
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
		want Level
	}{
		{s: "", want: LevelUnknown},
		{s: "nop", want: LevelNop},
		{s: "error", want: LevelError},
		{s: "warn", want: LevelWarn},
		{s: "info", want: LevelInfo},
		{s: "debug", want: LevelDebug},
		{s: "trace", want: LevelTrace},
		{s: "foo", want: LevelUnknown},
	}
	for _, tc := range testCases {
		if got := LevelFrom(tc.s); got != tc.want {
			t.Errorf("LevelFrom(\"%s\") is [%v], want [%v]", tc.s, got, tc.want)
		}
	}
}

func TestGologLevel(t *testing.T) {
	testCases := []struct {
		s    string
		want log.LogLevel
		err  error
	}{
		{s: "nop", want: 0, err: ecode.ErrLogLevel},
		{s: "error", want: log.LevelError, err: nil},
		{s: "warn", want: log.LevelWarn, err: nil},
		{s: "info", want: log.LevelInfo, err: nil},
		{s: "debug", want: log.LevelDebug, err: nil},
		{s: "trace", want: log.LevelDebug, err: nil},
	}
	for _, tc := range testCases {
		if lvl := LevelFrom(tc.s); lvl == LevelUnknown {
			t.Errorf("LevelFrom(\"%s\") is [%v]", tc.s, lvl)
		} else if got, err := lvl.GologLevel(); !errors.Is(err, tc.err) {
			t.Errorf("GologLevel(\"%v\") is [%v], want [%v]", tc.s, err, tc.err)
		} else if got != tc.want {
			t.Errorf("GologLevel(\"%v\") is [%v], want [%v]", tc.s, got, tc.want)
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
