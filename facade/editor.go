package facade

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/hymkor/go-multiline-ny"
)

func inputFromPipe(ui *rwi.RWI) (string, error) {
	b, err := io.ReadAll(ui.Reader())
	if err != nil {
		return "", errs.Wrap(err)
	}
	return string(b), nil
}

func editMessage(ctx context.Context, w io.Writer) (string, error) {
	var editor multiline.Editor
	editor.SetPrompt(func(w io.Writer, lnum int) (int, error) {
		return fmt.Fprintf(w, "%2d>", lnum+1)
	})
	fmt.Fprintln(w, "Input 'Ctrl+J' or 'Ctrl+Enter' to submit message")
	fmt.Fprintln(w, "Input 'Ctrl+D' with no chars to stop")
	lines, err := editor.Read(ctx)
	if err != nil {
		if errs.Is(err, io.EOF) {
			return "", nil
		}
		return "", errs.Wrap(err)
	}
	if len(lines) == 0 {
		return "", nil
	}
	return strings.Join(lines, "\n"), nil
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
