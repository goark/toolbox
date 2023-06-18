package feed

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/goark/errs"
	"github.com/goark/toolbox/ecode"
)

// Author is author data in Metadata.
type Author struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Image is image data in Metadata.
type Image struct {
	MimeType string     `json:"mime_type,omitempty"`
	Title    string     `json:"title,omitempty"`
	URL      string     `json:"url,omitempty"`
	FName    string     `json:"file_name,omitempty"`
	Taken    *time.Time `json:"taken,omitempty"`
}

// FileName returns file name from Image data.
func (i *Image) FileName() string {
	if i == nil {
		return ""
	}
	if len(i.FName) > 0 {
		return i.FName
	}
	u, err := url.Parse(i.URL)
	if err != nil {
		return ""
	}
	_, fname := path.Split(u.Path)
	return fname
}

// Item is item data in Metadata.
type Item struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Link        string     `json:"link,omitempty"`
	Authors     []*Author  `json:"authors,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
	Updated     *time.Time `json:"Updated,omitempty"`
	Images      []*Image   `json:"images,omitempty"`
}

// Metadata is metadata for feed.
type Metadata struct {
	FeedLink    string    `json:"feedLink,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Link        string    `json:"link,omitempty"`
	ID          string    `json:"id,omitempty"`
	Authors     []*Author `json:"authors,omitempty"`
	Items       []*Item   `json:"items,omitempty"`
}

// DecodeFromJSON decodes from JSON string to Metadata.
func DecodeFromJSON(r io.Reader) (*Metadata, error) {
	var data Metadata
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, errs.Wrap(err)
	}
	return &data, nil
}

// EncodeToJson encodes from Metadata to JSON string.
func (data *Metadata) EncodeToJson() (io.Reader, error) {
	if data == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, errs.Wrap(err)
	}
	return buf, nil
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
