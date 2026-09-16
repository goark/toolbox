package striptag

import (
	"errors"
	"html"
	"io"
	"strings"

	ghtml "golang.org/x/net/html"
)

// skipTags defines the HTML tags whose content should be skipped when stripping tags.
var skipTags = map[string]bool{
	"script":   true,
	"style":    true,
	"noscript": true,
	"iframe":   true,
	"object":   true,
	"embed":    true,
	"textarea": true,
	"title":    true,
}

// StripTags removes HTML tags from the input string while skipping
// the content of certain tags defined in skipTags.
func StripTags(s string) (string, error) {
	tokenizer := ghtml.NewTokenizer(strings.NewReader(s))
	var b strings.Builder
	b.Grow(len(s))
	skipDepth := 0

	for {
		tt := tokenizer.Next() // get the next token type
		switch tt {
		case ghtml.ErrorToken: // handle error token
			if errors.Is(tokenizer.Err(), io.EOF) {
				return b.String(), nil // return the accumulated text at the end of input
			}
			return "", tokenizer.Err()
		case ghtml.StartTagToken: // handle start tag token
			t := tokenizer.Token()
			if skipTags[t.Data] {
				skipDepth++ // start skipping content of this tag
			}
			if skipDepth == 0 {
				changeBrToNewline(t, &b)
			}
		case ghtml.SelfClosingTagToken: // handle self-closing tag token
			t := tokenizer.Token()
			if skipDepth == 0 {
				changeBrToNewline(t, &b)
			}
		case ghtml.EndTagToken: // handle end tag token
			t := tokenizer.Token()
			if skipDepth > 0 && skipTags[t.Data] {
				skipDepth-- // stop skipping one level of skipped content
			}
		case ghtml.TextToken: // handle text token
			if skipDepth == 0 {
				b.WriteString(html.UnescapeString(string(tokenizer.Text())))
			}
		}
	}
}

// changeBrToNewline handles the <br> tag by converting it to a newline character in the output.
func changeBrToNewline(t ghtml.Token, b *strings.Builder) {
	if t.Data == "br" {
		// Handle <br> as a line break.
		b.WriteByte('\n')
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
