package values

import (
	"encoding/json"
	"testing"
)

func TestDate(t *testing.T) {
	testCases := []struct {
		s     string
		isErr bool
	}{
		{s: "2023-02-22", isErr: false},
		{s: "2023", isErr: true},
		{s: "", isErr: false},
	}

	for _, tc := range testCases {
		s, err := DateFrom(tc.s, false)
		if (err != nil) != tc.isErr {
			t.Errorf("Is \"%v\" error ? %v, want %v", tc.s, err != nil, tc.isErr)
		}
		if err == nil {
			if s.String() != tc.s {
				t.Errorf("DateFrom(\"%v\") is \"%v\" , want \"%v\"", tc.s, s, tc.s)
			}
		}
	}
}

func TestDateJSON(t *testing.T) {
	testCases := []struct {
		s string
	}{
		{s: `{"date":"2023-02-22"}`},
		{s: `{"date":""}`},
	}

	for _, tc := range testCases {
		var data struct {
			Dt Date `json:"date"`
		}
		if err := json.Unmarshal([]byte(tc.s), &data); err != nil {
			t.Errorf("Unmarshal(\"%v\") is %v, want nil", tc.s, err)
		} else if b, err := json.Marshal(data); err != nil {
			t.Errorf("Marshal(\"%v\") is %v, want nil", tc.s, err)
		} else if string(b) != tc.s {
			t.Errorf("Unmarshal/Marshal is \"%v\", want \"%v\"", string(b), tc.s)
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
