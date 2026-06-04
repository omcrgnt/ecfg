package ecfgtool

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/omcrgnt/ecfg/internal/testdata"
)

func TestSetLeafValue_nilPointer(t *testing.T) {
	var p *testdata.Label
	err := setLeafValue(reflect.ValueOf(p), "x", false)
	if err == nil || !strings.Contains(err.Error(), "nil pointer") {
		t.Fatalf("got %v", err)
	}
}

func TestSetScalarValue_unsupportedKind(t *testing.T) {
	var r testdata.Ratio
	v := reflect.ValueOf(&r).Elem()
	err := setScalarValue(v, "1+2i")
	if err == nil || !strings.Contains(err.Error(), "unsupported scalar kind") {
		t.Fatalf("got %v", err)
	}
}

func TestSetScalarValue_duration(t *testing.T) {
	var d testdata.Timeout
	v := reflect.ValueOf(&d).Elem()
	if err := setScalarValue(v, "5s"); err != nil {
		t.Fatal(err)
	}
	if time.Duration(d) != 5*time.Second {
		t.Fatalf("got %v", d)
	}
}

func TestSetScalarValue_durationInvalid(t *testing.T) {
	var d testdata.Timeout
	v := reflect.ValueOf(&d).Elem()
	err := setScalarValue(v, "not-a-duration")
	if err == nil || !strings.Contains(err.Error(), "parse duration") {
		t.Fatalf("got %v", err)
	}
}

func TestSetScalarValue_floatInvalid(t *testing.T) {
	var s testdata.Score
	v := reflect.ValueOf(&s).Elem()
	err := setScalarValue(v, "not-float")
	if err == nil || !strings.Contains(err.Error(), "parse float") {
		t.Fatalf("got %v", err)
	}
}

func TestLookupEnv_missing(t *testing.T) {
	_, err := lookupEnv("ECFG_TEST_NONEXISTENT_KEY_XYZ")
	if !errors.Is(err, ErrMissingEnv) {
		t.Fatalf("got %v", err)
	}
}
