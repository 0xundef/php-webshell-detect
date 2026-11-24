package core

// AnalyzeItem is the analysis result for one script file.
type AnalyzeItem struct {
	Path       string
	CostTime   int64 // milliseconds
	HitTags    HIT_TAG
	SourceSink []string
	RegEval    []string
	DynamicFun []string
	Decode     []string
	LoopCount  int
	Confidence string
}

type AItems []AnalyzeItem

// String renders a human readable, line-oriented report.
func (a AnalyzeItem) String() string {
	var ret string
	ret = a.Path + "\n"
	for _, s := range a.SourceSink {
		ret += "type:normal " + s + "\n"
	}
	for _, s := range a.RegEval {
		ret += "type:regex eval" + s + "\n"
	}
	for _, s := range a.DynamicFun {
		ret += "type:dynamic fun" + s + "\n"
	}
	for _, s := range a.Decode {
		ret += "type:decode" + s + "\n"
	}
	return ret
}

// TagList returns the matched tags in a stable order.
func (a AnalyzeItem) TagList() []string {
	tags := make([]string, 0, 4)
	if a.HasTaintDim(SourceSink) {
		tags = append(tags, SourceSink.String())
	}
	if a.HasTaintDim(RegEval) {
		tags = append(tags, RegEval.String())
	}
	if a.HasTaintDim(SuspiciousDecode) {
		tags = append(tags, SuspiciousDecode.String())
	}
	if a.HasTaintDim(DynamicFun) {
		tags = append(tags, DynamicFun.String())
	}
	return tags
}

// Tags returns a comma separated tag list. It never has a leading comma.
func (a AnalyzeItem) Tags() string {
	var ret string
	for i, tag := range a.TagList() {
		if i > 0 {
			ret += ","
		}
		ret += tag
	}
	return ret
}

func (a AnalyzeItem) HasAny() bool {
	return a.HitTags != 0
}

func (a AnalyzeItem) HasTaintDim(tag HIT_TAG) bool {
	return a.HitTags&tag != 0
}

func (a AItems) Len() int {
	return len(a)
}

func (a AItems) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a AItems) Less(i, j int) bool {
	return a[i].CostTime > a[j].CostTime
}

// JsonItem is kept for backwards compatibility with batch consumers.
type JsonItem struct {
	Path       string
	HitData    string
	Confidence string
}
