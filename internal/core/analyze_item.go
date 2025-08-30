package core

// AnalyzeItem the checked result for each script file
type AnalyzeItem struct {
	Path       string
	CostTime   int64
	HitTags    HIT_TAG
	SourceSink []string
	RegEval    []string
	DynamicFun []string
	Decode     []string
	LoopCount  int
	Confidence string
}

type AItems []AnalyzeItem

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
func (a AnalyzeItem) Tags() string {
	var ret string
	if a.HasTaintDim(SourceSink) {
		ret = ret + "," + SourceSink.String()
	}
	if a.HasTaintDim(SuspiciousDecode) {
		ret = ret + "," + SuspiciousDecode.String()
	}
	if a.HasTaintDim(DynamicFun) {
		ret = ret + "," + DynamicFun.String()
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

//for output
type JsonItem struct {
	Path       string
	HitData    string
	Confidence string
}
