package striptag

import "testing"

func TestStripTags(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "br to newline",
			in:   "a<br>b",
			want: "a\nb",
		},
		{
			name: "br to newline",
			in:   "a<br/>b",
			want: "a\nb",
		},
		{
			name: "entities are unescaped",
			in:   "Tom &amp; Jerry",
			want: "Tom & Jerry",
		},
		{
			name: "skip nested tags content",
			in:   "<script>a<style>b</style>c</script>d",
			want: "d",
		},
		{
			name: "do not emit br inside skipped area",
			in:   "<script>before<br>after</script>x",
			want: "x",
		},
		{
			name: "malformed html still returns best effort",
			in:   "<div>abc",
			want: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StripTags(tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

/* Copyright 2026 Spiegel
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
