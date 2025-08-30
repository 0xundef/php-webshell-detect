package core

import (
	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"github.com/0xundef/php-webshell-detect/internal/core/context"
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/conf"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/version"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/php7"
	"github.com/0xundef/php-webshell-detect/frontend/php/version/scanner"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

var WORLD *World

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
		SourceSink: []string{},
		RegEval:    []string{},
		DynamicFun: []string{},
		Decode:     []string{},
		Confidence: "high",
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
	return result
}

func initialize(config DetectConf) *World {
	ir.PHP_CONTEXT = ir.InitPHPContext()
	if common.Log() == nil {
		common.InitDefaultLogger()
	}
	hierarchy := lang.CreateClassHierarchy()
	typeSystem := &lang.PHPTypeSystem{
		Hierarchy: hierarchy,
	}
	var selector context.ContextSelector
	heap := heap.CreateHeapModel()
	if config.ContextPolicy == "ci" {
		selector = context.CreateCISelector()
	} else if strings.Contains(config.ContextPolicy, "k-callsite") {
		items := strings.Split(config.ContextPolicy, "-")
		if len(items) == 3 {
			limit, _ := strconv.Atoi(items[2])
			selector = context.CreateKCallSiteSelector(limit)
		}
	}
	solver := CreateSolver(heap, selector, context.CreateCSManager())

	w := &World{
		TypeSystem: typeSystem,
		Hierarchy:  hierarchy,
		MainMethod: nil,
		RootCtx:    context.CreateEmptyContext(),
		S:          solver,

		LoopLimit:  config.LoopLimit,
		TargetFile: config.Path,
		Debug:      config.Debug,
		MockPath:   config.ConfigDir + "/mock.php",
		IRS:        map[string]ir.IR{},
	}
	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	solver.world = w
	return w
}

func Bang(config DetectConf) *World {
	w := initialize(config)
	w.buildAllClasses()
	return w
}

func (w *World) buildAllClasses() {
	config := conf.Config{
		Version: &version.Version{
			Major: 7,
			Minor: 4,
		},
	}
	path := w.TargetFile
	content, err := os.ReadFile(path)
	if err != nil {
		panic(path + " dose not exist")
	}
	mock, err := os.ReadFile(w.MockPath)
	if err != nil {
		panic("config/mock.php read failed")
	}
	content = append(content, mock...)
	lexer := scanner.NewLexer(content, config)
	php7parser := php7.NewParser(lexer, config)
	php7parser.Parse()
	root := php7parser.GetRootNode()

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
}
