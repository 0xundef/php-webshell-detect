package core

type HIT_TAG int

// Hit tags reported by the detection pass.
const (
	// All matches every tag (used by batch output filters).
	All HIT_TAG = -1
	// SourceSink: tainted value reaches a sensitive parameter.
	SourceSink HIT_TAG = 1
	// DynamicFun: call through a computed function name.
	DynamicFun HIT_TAG = 1 << 1
	// SuspiciousDecode: decode/transform chain reaches a sink.
	SuspiciousDecode HIT_TAG = 1 << 2
	// RegEval: regex with the eval modifier (e.g. /x/e).
	RegEval HIT_TAG = 1 << 3
)

// Confidence levels reported alongside the hit tags.
const (
	ConfidenceHigh   = "high"
	ConfidenceMedium = "medium"
	ConfidenceNone   = "none"
)

func (h HIT_TAG) String() string {
	switch h {
	case SourceSink:
		return "SourceSink"
	case DynamicFun:
		return "DynamicFun"
	case SuspiciousDecode:
		return "SuspiciousDecode"
	case RegEval:
		return "RegEval"
	}
	return ""
}
