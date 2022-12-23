package schema

import "testing"

func TestHdnsTimeUnmarshalJSON(t *testing.T) {
	passUnmarshalTime(t, "2020-04-07 01:24:37 +0000 UTC")
	passUnmarshalTime(t, "2020-04-07 01:56:03.196438163 +0000 UTC m=+755.322810452")
	passUnmarshalTime(t, "2022-12-13 01:37:45.814 +0000 UTC")
	passUnmarshalTime(t, "")

	failUnMarshalTime(t, "random text")
	failUnMarshalTime(t, "01-20-2020 02:32:48")
}

func passUnmarshalTime(t *testing.T, s string) {
	ti := &HdnsTime{}
	if err := ti.UnmarshalJSON([]byte(s)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func failUnMarshalTime(t *testing.T, s string) {
	ti := &HdnsTime{}
	if err := ti.UnmarshalJSON([]byte(s)); err == nil {
		t.Errorf("missing expected error")
	}
}
