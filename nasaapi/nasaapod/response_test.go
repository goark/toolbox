package nasaapod

import (
	"errors"
	"strings"
	"testing"

	"github.com/goark/toolbox/nasaapi"
)

func TestDecodeNormalArray(t *testing.T) {
	raw := `[{"date":"2026-09-19","title":"sample"}]`
	res, err := decode(strings.NewReader(raw), false)
	if err != nil {
		t.Fatalf("decode() error = %v, want nil", err)
	}
	if len(res) != 1 {
		t.Fatalf("len(response) = %d, want 1", len(res))
	}
	if res[0].Title != "sample" {
		t.Fatalf("title = %q, want %q", res[0].Title, "sample")
	}
}

func TestDecodeNormalSingle(t *testing.T) {
	raw := `{"date":"2026-09-19","title":"single"}`
	res, err := decode(strings.NewReader(raw), true)
	if err != nil {
		t.Fatalf("decode() error = %v, want nil", err)
	}
	if len(res) != 1 {
		t.Fatalf("len(response) = %d, want 1", len(res))
	}
	if res[0].Title != "single" {
		t.Fatalf("title = %q, want %q", res[0].Title, "single")
	}
}

func TestDecodeErrorByStatus(t *testing.T) {
	raw := `{"code":"","message":"Invalid parameter(s)","data":{"status":400}}`
	_, err := decode(strings.NewReader(raw), false)
	if !errors.Is(err, nasaapi.ErrAPODAPIResponse) {
		t.Fatalf("decode() error = %v, want ErrAPODAPIResponse", err)
	}
}

func TestDecodeErrorByCodeWithoutStatus(t *testing.T) {
	raw := `{"code":"rest_invalid_param","message":"Invalid parameter(s)","data":{"status":0}}`
	_, err := decode(strings.NewReader(raw), false)
	if !errors.Is(err, nasaapi.ErrAPODAPIResponse) {
		t.Fatalf("decode() error = %v, want ErrAPODAPIResponse", err)
	}
}

func TestDecodeNonErrorObjectDoesNotUseAPODError(t *testing.T) {
	raw := `{"code":"","message":"notice","data":{"status":0}}`
	_, err := decode(strings.NewReader(raw), false)
	if err == nil {
		t.Fatalf("decode() error = nil, want non-nil")
	}
	if errors.Is(err, nasaapi.ErrAPODAPIResponse) {
		t.Fatalf("decode() error = %v, should not be ErrAPODAPIResponse", err)
	}
}
