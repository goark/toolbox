package feed_test

import (
	"testing"

	"github.com/goark/toolbox/webpage/feed"
)

func TestImageFileName(t *testing.T) {
	testCases := []struct {
		i *feed.Image
		f string
	}{
		{i: &feed.Image{}, f: ""},
		{i: &feed.Image{FName: "foo"}, f: "foo"},
		{i: &feed.Image{FName: "foo", URL: "http://foo/bar/image.jpg"}, f: "foo"},
		{i: &feed.Image{FName: "", URL: "http://foo/bar/image.jpg"}, f: "image.jpg"},
	}

	for _, tc := range testCases {
		f := tc.i.FileName()
		if f != tc.f {
			t.Errorf("Image.FileName() = \"%v\", want \"%v\".", f, tc.f)
		}
	}
}

/* Copyright 2022 Spiegel
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
