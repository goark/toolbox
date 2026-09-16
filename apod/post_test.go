package apod

import (
	"strings"
	"testing"

	"github.com/goark/toolbox/nasaapi/nasaapod"
	"github.com/goark/toolbox/values"
)

func dateMust(s string) values.Date {
	d, err := values.DateFrom(s, false)
	if err != nil {
		panic(err)
	}
	return d
}

func TestMakeMessageNil(t *testing.T) {
	msg, err := MakeMessage(nil)
	if err != nil {
		t.Fatalf("MakeMessage(nil) error = %v, want nil", err)
	}
	if msg != "" {
		t.Fatalf("MakeMessage(nil) msg = %q, want empty", msg)
	}
}

func TestMakeMessageStripCreditHTML(t *testing.T) {
	// #nosec G101 -- test fixture intentionally includes an HTML credit snippet.
	res := &nasaapod.Response{
		Date:      dateMust("2026-09-16"),
		Title:     "Test APOD",
		MediaType: nasaapod.MediaImage,
		Credit:    `<a href="https://example.com">Alice &amp; Bob</a>`,
		Permalink: "https://science.nasa.gov/image-article/example/",
	}

	msg, err := MakeMessage(res)
	if err != nil {
		t.Fatalf("MakeMessage() error = %v, want nil", err)
	}
	if strings.Contains(msg, "<a") {
		t.Fatalf("message contains raw HTML: %q", msg)
	}
	if !strings.Contains(msg, "Image Credit: Alice & Bob") {
		t.Fatalf("message credit line mismatch: %q", msg)
	}
	if !strings.Contains(msg, "Web page: https://science.nasa.gov/image-article/example/") {
		t.Fatalf("message web page line missing: %q", msg)
	}
	if strings.Contains(msg, "Content:") {
		t.Fatalf("image media should not include content line: %q", msg)
	}
}

func TestMakeMessageVideoIncludesContent(t *testing.T) {
	res := &nasaapod.Response{
		Date:      dateMust("2026-09-16"),
		Title:     "Video APOD",
		MediaType: nasaapod.MediaVideo,
		Url:       "https://example.com/video",
		Permalink: "https://science.nasa.gov/image-article/example-video/",
	}

	msg, err := MakeMessage(res)
	if err != nil {
		t.Fatalf("MakeMessage() error = %v, want nil", err)
	}
	if !strings.Contains(msg, "Content: https://example.com/video") {
		t.Fatalf("video media should include content line: %q", msg)
	}
}
