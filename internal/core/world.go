package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/conf"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/version"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/php7"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/scanner"
	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"github.com/0xundef/php-webshell-detect/internal/core/context"
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"github.com/sirupsen/logrus"
)

type World struct {
	TypeSystem lang.TypeSystem
	Hierarchy  *lang.PHPClassHierarchy
	MainMethod *lang.PHPMethod
	RootCtx    context.Context
	S          *Solver

	LoopLimit  int
	TargetFile string
	IRS        map[string]ir.IR
	Debug      bool
	NeedMock   bool
	MockPath   string
}

func (w *World) GetIR(method *lang.PHPMethod) ir.IR {
	if ir, ok := w.IRS[method.GetSignature()]; ok {
		return ir
	}
	if method.Source == nil {
		return nil
	}
	this := expr.Var{
		Name: "$this",
		Type: nil,
	}
	visitor := &ir.AstVisitor{
		Ctx: ir.VisitContext{},
		Ir: &ir.DefaultIR{
			Method:  method,
			Stmts:   []stmt.Stmt{},
			Params:  []expr.Var{},
			RetVars: []expr.Var{},
			This:    this,
		},
		Hierarchy: w.Hierarchy,
		VM:        w.S.VM,
		Debug:     w.Debug,
	}
	method.Source.Accept(visitor)

	ir := visitor.Ir.(*ir.DefaultIR)
	for _, name := range method.ParamNames {
		v := expr.Var{
			Name: name,
			Type: nil,
		}
		ir.Params = append(ir.Params, v)
	}
	w.IRS[method.GetSignature()] = ir
	return ir
}

// Analyze runs the solver and collects findings. It may panic on internal
// invariant violations; callers that handle untrusted input should use
// AnalyzeSafe instead.
func (w *World) Analyze() AnalyzeItem {
	w.S.Solve()
	ret := w.Collect()
	ret.LoopCount = w.S.LoopCount

	if w.Debug {
		common.Log().Info("================ Print IR begin ================")
		w.PrintIRS()
		common.Log().Info("================ Print IR end   ================")
		common.Log().Info("================ Print var pointer begin ================")
		w.PrintVarDetail()
		common.Log().Info("================ Print var pointer end   ================")
		common.Log().Info("================ Print CallEdges begin   ================")
		w.PrintAllCallEdges()
		common.Log().Info("================ Print CallEdges end     ================")

		if ret.LoopCount > w.LoopLimit {
			common.Log().Info("has exceed loop limit")
		}
	}
	return ret
}
func (w *World) PrintAllCallEdges() {
	for _, edge := range w.S.CallEdges {
		for _, name := range edge.Callee.Method.ParamNames {
			v := expr.Var{
				Name:   name,
				Type:   nil,
				Global: false,
			}
			csArg := w.S.CSM.GetCSVar(edge.Callee.Ctx, v)
			common.Log().Info(csArg.String() + ":" + strconv.Itoa(csArg.MergeCount) + "--->" + "[" + csArg.PointsToSet.Sprint() + "]")
		}
	}
}

func (w *World) PrintIRS() {
	for _, ir := range w.IRS {
		common.Log().Info(ir.String())
	}
}
func (w *World) PrintVarDetail() {
	for _, pointer := range w.S.CSM.PrtManager.Vars {
		common.Log().Info(pointer.String() + "--->" + pointer.PointsToSet.Sprint())
	}
}
func (w *World) Collect() AnalyzeItem {
	result := AnalyzeItem{
		Path:       w.TargetFile,
		SourceSink: []string{},
		RegEval:    []string{},
		DynamicFun: []string{},
		Decode:     []string{},
	}
	for _, edge := range w.S.CallEdges {
		for _, name := range edge.Callee.Method.ParamNames {
			v := expr.Var{
				Name:   name,
				Type:   nil,
				Global: false,
			}
			csArg := w.S.CSM.GetCSVar(edge.Callee.Ctx, v)
			if strings.Contains(csArg.GetVar().Name, "$sink") && strings.Contains(csArg.PointsToSet.Sprint(), "taint") {
				if !result.HasTaintDim(SourceSink) {
					result.HitTags = result.HitTags | SourceSink
				}
				result.SourceSink = append(result.SourceSink, "	"+csArg.String()+"--->"+"["+csArg.PointsToSet.Sprint()+"]")
			}

			if strings.Contains(csArg.GetVar().Name, "$pattern") && (csArg.PointsToSet.SprintIfHasEndWith("/e") != "" || csArg.PointsToSet.SprintIfHasEndWith("|e") != "") {
				if !result.HasTaintDim(RegEval) {
					result.HitTags = result.HitTags | RegEval
				}
				result.RegEval = append(result.RegEval, "	"+csArg.String()+"--->"+"["+csArg.PointsToSet.Sprint()+"]")
			}

			if strings.Contains(csArg.GetVar().Name, "$options") {
				record := csArg.PointsToSet.SprintIfHasEqual("e")
				if record != "" {
					if !result.HasTaintDim(RegEval) {
						result.HitTags = result.HitTags | RegEval
					}
					result.RegEval = append(result.RegEval, "	"+csArg.String()+"--->"+"["+record+"]")
				}
			}
			if strings.Contains(csArg.GetVar().Name, "$sink") && (strings.Contains(csArg.PointsToSet.Sprint(), "decode") || strings.Contains(csArg.PointsToSet.Sprint(), "uncompress")) {
				if !result.HasTaintDim(SuspiciousDecode) {
					result.HitTags = result.HitTags | SuspiciousDecode
				}
				result.Decode = append(result.Decode, "	"+csArg.String()+"--->"+"["+csArg.PointsToSet.Sprint()+"]")
			}
		}

		if invoke, ok := edge.CallSite.CallSite.Rvalue.(*expr.InvokeStaticExp); ok {
			if invoke.FunVar.Name != "" {
				a := w.S.CSM.GetCSVar(edge.CallSite.Ctx, invoke.FunVar)
				if a.GetMergeCount() > 1 || strings.Contains(a.PointsToSet.Sprint(), "replace") ||
					strings.Contains(a.PointsToSet.Sprint(), "decode") {
					if !result.HasTaintDim(DynamicFun) {
						result.HitTags = result.HitTags | DynamicFun
					}
					result.DynamicFun = append(result.DynamicFun, "	"+edge.CallSite.String()+":"+strconv.Itoa(a.GetMergeCount()))
				}
				if strings.Contains(a.PointsToSet.Sprint(), "taint") {
					if !result.HasTaintDim(DynamicFun) {
						result.HitTags = result.HitTags | DynamicFun
					}
					result.DynamicFun = append(result.DynamicFun, "	"+edge.CallSite.String()+":"+strconv.Itoa(a.GetMergeCount()))
				}
				if a.GetMergeCount() >= 1 {
					for _, name := range edge.Callee.Method.ParamNames {
						v := expr.Var{
							Name:   name,
							Type:   nil,
							Global: false,
						}
						csArg := w.S.CSM.GetCSVar(edge.Callee.Ctx, v)
						if strings.Contains(csArg.PointsToSet.Sprint(), "taint") {
							if !result.HasTaintDim(DynamicFun) {
								result.HitTags = result.HitTags | DynamicFun
							}
							result.DynamicFun = append(result.DynamicFun, "	"+edge.CallSite.String()+":"+strconv.Itoa(a.GetMergeCount()))
						}
					}
				}
			}
		}
	}
	result.Confidence = confidenceFor(result.HitTags)
	return result
}

// confidenceFor maps the coarse hit tags to a confidence level. A direct
// source-to-sink or a regex-eval hit is high confidence; obfuscation-only
// evidence is medium; no hit at all is "none".
func confidenceFor(tags HIT_TAG) string {
	switch {
	case tags&(SourceSink|RegEval) != 0:
		return ConfidenceHigh
	case tags&(SuspiciousDecode|DynamicFun) != 0:
		return ConfidenceMedium
	default:
		return ConfidenceNone
	}
}

// phpContextOnce guards the process-wide PHP builtin list. It is immutable
// after the first initialization, so concurrent analyses can share it safely.
var phpContextOnce sync.Once

// DefaultLoopLimit bounds solver iterations when the config omits a limit.
const DefaultLoopLimit = 100000

func initialize(config DetectConf) (*World, error) {
	if config.Path == "" {
		return nil, fmt.Errorf("no target file supplied")
	}
	if config.LoopLimit <= 0 {
		config.LoopLimit = DefaultLoopLimit
	}
	if config.ConfigDir == "" {
		config.ConfigDir = "config"
	}

	selector, err := context.NewSelector(config.ContextPolicy)
	if err != nil {
		return nil, err
	}

	phpContextOnce.Do(func() {
		ir.PHP_CONTEXT = ir.InitPHPContext()
	})
	if common.Log() == nil {
		common.InitDefaultLogger()
	}
	hierarchy := lang.CreateClassHierarchy()
	typeSystem := &lang.PHPTypeSystem{
		Hierarchy: hierarchy,
	}
	hm := heap.CreateHeapModel()
	solver := CreateSolver(hm, selector, context.CreateCSManager())

	w := &World{
		TypeSystem: typeSystem,
		Hierarchy:  hierarchy,
		MainMethod: nil,
		RootCtx:    context.CreateEmptyContext(),
		S:          solver,

		LoopLimit:  config.LoopLimit,
		TargetFile: config.Path,
		Debug:      config.Debug,
		MockPath:   filepath.Join(config.ConfigDir, "mock.php"),
		IRS:        map[string]ir.IR{},
	}
	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	solver.world = w
	return w, nil
}

// New parses the target file together with the mock definitions and returns a
// world ready for Analyze. User-facing failures are returned as errors; a
// panic inside the frontend is converted into an error as well.
func New(config DetectConf) (w *World, err error) {
	defer func() {
		if r := recover(); r != nil {
			w = nil
			err = fmt.Errorf("initialize analysis: %v", r)
		}
	}()

	w, err = initialize(config)
	if err != nil {
		return nil, err
	}
	if err = w.buildAllClasses(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *World) buildAllClasses() error {
	config := conf.Config{
		Version: &version.Version{
			Major: 7,
			Minor: 4,
		},
	}
	content, err := os.ReadFile(w.TargetFile)
	if err != nil {
		return fmt.Errorf("read target file %q: %w", w.TargetFile, err)
	}
	mock, err := os.ReadFile(w.MockPath)
	if err != nil {
		return fmt.Errorf("read mock file %q: %w", w.MockPath, err)
	}
	content = append(content, mock...)
	lexer := scanner.NewLexer(content, config)
	php7parser := php7.NewParser(lexer, config)
	php7parser.Parse()
	root := php7parser.GetRootNode()
	if root == nil {
		return fmt.Errorf("parse %q: empty syntax tree", w.TargetFile)
	}

	entryClass := lang.CreatePHPClass("EntryMain", "main")
	mainMethod := lang.CreatePHPMethod(entryClass, "main", nil, nil, nil)
	mainMethod.Source = root
	w.MainMethod = mainMethod
	entryClass.AddMethod(mainMethod)
	w.Hierarchy.AddClass(entryClass)

	//build all Classes
	converter := &ir.AstVisitor{
		Ctx:        ir.VisitContext{},
		Ir:         nil,
		Hierarchy:  w.Hierarchy,
		BuildClass: true,
		VM:         w.S.VM,
		Debug:      w.Debug,
	}
	root.Accept(converter)
	return nil
}

// AnalyzeSafe runs the solver and converts internal panics into errors so that
// callers do not have to recover themselves.
func (w *World) AnalyzeSafe() (item AnalyzeItem, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("analysis failed: %v", r)
		}
	}()
	return w.Analyze(), nil
}
