package apod

import (
	"fmt"
	"strings"

	"github.com/goark/toolbox/nasaapi/nasaapod"
	"github.com/goark/toolbox/striptag"
)

// MakeMessage creates a message string from the given APOD response data.
// It includes the hash tag, title, credit, web page, and content URL.
func MakeMessage(data *nasaapod.Response) string {
	if data == nil {
		return ""
	}
	bld := strings.Builder{}

	// hash tag
	bld.WriteString("#apod ")
	bld.WriteString(data.Date.String())
	if len(data.MediaType) > 0 && data.MediaType != "image" {
		fmt.Fprintf(&bld, " (%s)", data.MediaType)
	}
	bld.WriteString("\n")
	// title
	if len(data.Title) > 0 {
		fmt.Fprintln(&bld, data.Title)
	}
	// credit
	credit := data.Credit
	if len(credit) == 0 {
		credit = data.Copyright
	}
	if len(credit) > 0 {
		txt, err := striptag.StripTags(credit) // remove HTML tags from credit
		if err == nil {
			credit = txt
		}
		fmt.Fprintln(&bld, "Image Credit:", credit)
	}
	// Web page
	fmt.Fprintln(&bld, "Web page:", data.WebPage())
	// content URL
	if data.MediaType != nasaapod.MediaImage && len(data.Url) > 0 {
		fmt.Fprintln(&bld, "Content:", data.Url)
	}
	return bld.String()
}

/* Copyright 2023-2026 Spiegel
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
