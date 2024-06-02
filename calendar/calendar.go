package calendar

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/goark/errs"
	"github.com/goark/koyomi"
	"github.com/goark/koyomi/value"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/tempdir"
)

const defaultTemplate = `{{ range . }}{{ .Date }} {{ .Title }}
{{ end }}`

// Config type is configurations for calendar package.
type Config struct {
	start, end         value.DateJp
	holiday, ephemeris bool
	tempDir            *tempdir.TempDir
	templateFile       string
}

// NewConfig function creates new Config instance.
func NewConfig(start, end string, holiday, ephemeris bool, tempDir *tempdir.TempDir, templateFile string) (*Config, error) {
	startdate := value.NewDate(time.Time{})
	enddate := value.NewDate(time.Time{})
	if len(start) == 0 && len(end) == 0 {
		tm := time.Now()
		startdate = value.NewDateEra(value.EraUnknown, tm.Year(), tm.Month(), 1)
		enddate = value.NewDateEra(value.EraUnknown, tm.Year(), tm.Month()+1, 0)
	} else {
		if len(start) > 0 {
			dt, err := value.DateFrom(start)
			if err != nil {
				return nil, errs.Wrap(err, errs.WithContext("start", start), errs.WithContext("end", end), errs.WithContext("holiday", holiday), errs.WithContext("ephemeris", ephemeris), errs.WithContext("templateFile", templateFile))
			}
			startdate = dt
		}
		if len(start) > 0 {
			dt, err := value.DateFrom(end)
			if err != nil {
				return nil, errs.Wrap(err, errs.WithContext("start", start), errs.WithContext("end", end), errs.WithContext("holiday", holiday), errs.WithContext("ephemeris", ephemeris), errs.WithContext("templateFile", templateFile))
			}
			enddate = dt
		}
	}
	return &Config{
		start:        startdate,
		end:          enddate,
		holiday:      holiday,
		ephemeris:    ephemeris,
		tempDir:      tempDir,
		templateFile: templateFile,
	}, nil
}

func (cal *Config) OutputEvent(w io.Writer) error {
	if cal == nil {
		return errs.Wrap(ecode.ErrNullPointer)
	}
	evs, err := cal.GetEvents()
	if err != nil {
		return errs.Wrap(err)
	}

	var tpl *template.Template
	if len(cal.templateFile) > 0 {
		t, err := template.New(filepath.Base(cal.templateFile)).ParseFiles(cal.templateFile)
		if err != nil {
			return errs.Wrap(err, errs.WithContext("templateFile", cal.templateFile))
		}
		tpl = t
	} else {
		t, err := template.New("").Parse(defaultTemplate)
		if err != nil {
			return errs.Wrap(err, errs.WithContext("defaultTemplate", defaultTemplate))
		}
		tpl = t
	}
	if err := tpl.Execute(w, evs); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func OutputEventSimple(evs []koyomi.Event) string {
	b := strings.Builder{}
	for _, ev := range evs {
		b.WriteString(fmt.Sprintln(ev.Date, ev.Title))
	}
	return b.String()
}

// GetEvents method returns astronomical events.
func (cal *Config) GetEvents() ([]koyomi.Event, error) {
	if cal == nil {
		return nil, errs.Wrap(ecode.ErrNullPointer)
	}
	ids := []koyomi.CalendarID{}
	if cal.holiday {
		ids = append(ids, koyomi.Holiday)
	}
	if cal.ephemeris {
		ids = append(ids, koyomi.MoonPhase, koyomi.SolarTerm, koyomi.Eclipse, koyomi.Planet)
	}
	if len(ids) == 0 {
		return []koyomi.Event{}, nil
	}
	k, err := koyomi.NewSource(
		koyomi.WithCalendarID(ids...),
		koyomi.WithStartDate(cal.start),
		koyomi.WithEndDate(cal.end),
		koyomi.WithTempDir(cal.tempDir.Path()),
	).Get()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return k.Events(), nil
}

/* Copyright 2024 Spiegel
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
