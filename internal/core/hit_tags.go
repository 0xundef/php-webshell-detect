package core

type HIT_TAG int

//SourceSink:normal type
//RegEval:regex match eval option
//DynamicFun:call suspicious dynamic function
//SuspiciousDecode:pass suspicious arguments to sink function
const (
	All              HIT_TAG = -1
	SourceSink       HIT_TAG = 1
	DynamicFun       HIT_TAG = 1 << 1
	SuspiciousDecode HIT_TAG = 1 << 2
	RegEval          HIT_TAG = 1 << 3
)

func (h HIT_TAG) String() string {
	switch h {
	case 1:
		return "SourceSink"
	case 2:
		return "DynamicFun"
	case 4:
		return "SuspiciousDecode"
	case 8:
		return "RegEval"
	}
	return ""
}
