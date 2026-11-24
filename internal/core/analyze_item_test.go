package core

import (
	"reflect"
	"testing"
)

func TestTagsOrderAndFormat(t *testing.T) {
	item := AnalyzeItem{HitTags: SourceSink | RegEval | SuspiciousDecode | DynamicFun}
	want := []string{"SourceSink", "RegEval", "SuspiciousDecode", "DynamicFun"}
	if got := item.TagList(); !reflect.DeepEqual(got, want) {
		t.Fatalf("TagList = %v, want %v", got, want)
	}
	if got, want := item.Tags(), "SourceSink,RegEval,SuspiciousDecode,DynamicFun"; got != want {
		t.Fatalf("Tags = %q, want %q", got, want)
	}
}

func TestTagsEmpty(t *testing.T) {
	item := AnalyzeItem{}
	if got := item.Tags(); got != "" {
		t.Fatalf("Tags = %q, want empty", got)
	}
	if got := item.TagList(); len(got) != 0 {
		t.Fatalf("TagList = %v, want empty", got)
	}
}

func TestConfidenceFor(t *testing.T) {
	cases := []struct {
		tags HIT_TAG
		want string
	}{
		{0, ConfidenceNone},
		{SourceSink, ConfidenceHigh},
		{RegEval, ConfidenceHigh},
		{SourceSink | SuspiciousDecode, ConfidenceHigh},
		{DynamicFun, ConfidenceMedium},
		{SuspiciousDecode, ConfidenceMedium},
		{DynamicFun | RegEval, ConfidenceHigh},
	}
	for _, tc := range cases {
		if got := confidenceFor(tc.tags); got != tc.want {
			t.Errorf("confidenceFor(%d) = %q, want %q", tc.tags, got, tc.want)
		}
	}
}
