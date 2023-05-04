package facade

import (
	"bytes"
	"testing"

	"github.com/goark/gocli/exitcode"
	"github.com/goark/gocli/rwi"
)

func TestNoCommand(t *testing.T) {
	result := "Error: no command\n"

	outBuf := new(bytes.Buffer)
	outErrBuf := new(bytes.Buffer)
	ui := rwi.New(rwi.WithWriter(outBuf), rwi.WithErrorWriter(outErrBuf))
	args := []string{}

	exit := Execute(ui, args)
	if exit != exitcode.Abnormal {
		t.Errorf("Execute() = \"%v\", want \"%v\".", exit, exitcode.Abnormal)
	}
	str := outBuf.String()
	if str != "" {
		t.Errorf("Execute() = \"%v\", want \"%v\".", str, "")
	}
	str = outErrBuf.String()
	if str != result {
		t.Errorf("Execute() = \"%v\", want \"%v\".", str, result)
	}
}

func TestNoCommandDebug(t *testing.T) {
	result := `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"no command"},"Context":{"function":"github.com/goark/toolbox/facade.newRootCmd.func1"}}` + "\n"

	outBuf := new(bytes.Buffer)
	outErrBuf := new(bytes.Buffer)
	ui := rwi.New(rwi.WithWriter(outBuf), rwi.WithErrorWriter(outErrBuf))
	args := []string{"--debug"}

	exit := Execute(ui, args)
	if exit != exitcode.Normal {
		t.Errorf("Execute() = \"%v\", want \"%v\".", exit, exitcode.Normal)
	}
	str := outBuf.String()
	if str != result {
		t.Errorf("Execute() = \"%v\", want \"%v\".", str, result)
	}
	str = outErrBuf.String()
	if str != "" {
		t.Errorf("Execute() = \"%v\", want \"%v\".", str, "")
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
