package core

import (
	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"github.com/0xundef/php-webshell-detect/internal/core/common/queue"
	"github.com/0xundef/php-webshell-detect/internal/core/context"
	"github.com/0xundef/php-webshell-detect/internal/core/heap"
	"github.com/0xundef/php-webshell-detect/internal/core/ir"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
)

// Solver the core engine
type Solver struct {
	WList     WorkList
	PFG       context.PointerGraph
	Selector  context.ContextSelector
	CSM       *context.CSManager
	HM        heap.HeapModel
	VM        *ir.VarManager
	LoopCount int
	world     *World

	CallEdges []context.CallEdge
}

func CreateSolver(model heap.HeapModel, selector context.ContextSelector, manager *context.CSManager) *Solver {
	workList := WorkList{
		PointerEntries: map[string]*context.PointsToSet{},
		PointerMap:     map[string]context.Pointer{},
		CallEdges:      queue.Queue{},
	}
	return &Solver{
		Selector: selector,
		CSM:      manager,
		HM:       model,
		WList:    workList,
		PFG:      context.PointerGraph{Pointers: map[string]context.Pointer{}},
		VM:       ir.CreateVarManager(),
	}
}

// AddEntryPoint set the entry method to analyze, usually the main method which is a mocked method for the top level source code in a script file
func (s *Solver) AddEntryPoint(entryMethod *lang.PHPMethod) {
	entryCtx := s.Selector.GetEmptyContext()
	csMethod := s.CSM.GetCSMethod(entryCtx, entryMethod)
	s.AddCSMethod(csMethod)
}

// Solve main algorithm implement P/Taint
func (s *Solver) Solve() {
	s.AddEntryPoint(s.world.MainMethod)
	var entry Entry
	entry = s.WList.PollEntry()
	for entry != nil {
		s.LoopCount++
		if s.LoopCount > s.world.LoopLimit {
			return
		}

		if pEntry, ok := entry.(*PointerEntry); ok { // propagate the diff new values to related pointers
			p := pEntry.GetPointer()
			pts := pEntry.GetPointsToSet()
			diff := s.propagate(p, pts)
			if !diff.IsEmpty() {
				if cvar, ok := p.(*context.CSVarPointer); ok {
					s.processInstanceStore(cvar, diff)
					s.processInstanceLoad(cvar, diff)
					s.processArrayStore(cvar, diff)
					s.processArrayLoad(cvar, diff)
					s.processCall(cvar, diff) //will generate more CallEdge
					s.processMerge(cvar, diff)
				}
			}
		} else if callEntry, ok := entry.(*CallEdgeEntry); ok { // extend new method or function calls, will involve
			s.processCallEdge(callEntry.CallEdge)
		}
		entry = s.WList.PollEntry()
	}
}

// o.f = x, o has diff obj then update pfg x->o_1.f, x->o_2.f
func (s *Solver) processInstanceStore(cvar *context.CSVarPointer, diff *context.PointsToSet) {
	var field *lang.PHPClassField
	for _, storeStmt := range cvar.GetStoreField() {
		sVar := storeStmt.Rvalue
		sP := s.CSM.GetCSVar(cvar.Context, sVar)
		for _, csObj := range diff.Set0.ToSlice() {
			fName := storeStmt.Lvalue.FieldRef.Name
			newObj := csObj.Obj.(*heap.NewObj)
			if class, ok := newObj.GetObjType().(*lang.ClassType); ok {
				field = s.world.Hierarchy.LookupField(class.Class, fName)
				//if field != nil {
				tP := s.CSM.GetInstanceField(csObj, field)
				s.AddPFGEdge(sP, tP, 1)
				//}
			}
		}
	}
}

func (s *Solver) processInstanceLoad(cvar *context.CSVarPointer, diff *context.PointsToSet) {
	var field *lang.PHPClassField
	for _, loadStmt := range cvar.GetLoadFields() {
		toVar := loadStmt.Lvalue
		tP := s.CSM.GetCSVar(cvar.Context, toVar)
		for _, csObj := range diff.Set0.ToSlice() {
			fName := loadStmt.Rvalue.FieldRef.Name
			newObj := csObj.Obj.(*heap.NewObj)
			if class, ok := newObj.GetObjType().(*lang.ClassType); ok {
				field = s.world.Hierarchy.LookupField(class.Class, fName)
				if field != nil {
					sP := s.CSM.GetInstanceField(csObj, field)
					s.AddPFGEdge(sP, tP, 1)
				}
			}
		}
	}
}

func (s *Solver) processArrayStore(cvar *context.CSVarPointer, diff *context.PointsToSet) {
	for _, arrayStmt := range cvar.GetStoreArrays() {
		sVar := arrayStmt.Rvalue
		sP := s.CSM.GetCSVar(cvar.Context, sVar)
		for _, csObj := range diff.Set0.ToSlice() {
			tP := s.CSM.GetArrayPointer(csObj)
			s.AddPFGEdge(sP, tP, 1)
		}
	}
}

// a.age , a.Name; a is cvar as based var, diff is new Obj that cvar pointed to, iterate all the Field in cvar and create new instance.Field point
func (s *Solver) processArrayLoad(cvar *context.CSVarPointer, diff *context.PointsToSet) {
	for _, arrayStmt := range cvar.GetLoadArrays() {
		toVar := arrayStmt.Stmt.(*stmt.LoadArrayStmt).Lvalue
		tP := s.CSM.GetCSVar(arrayStmt.Ctx, toVar)
		for _, csObj := range diff.Set0.ToSlice() {
			sP := s.CSM.GetArrayPointer(csObj)
			s.AddPFGEdge(sP, tP, 1)
		}
	}
}

func (s *Solver) processMerge(cvar *context.CSVarPointer, diff *context.PointsToSet) {
	for _, mergeStmt := range cvar.GetMerge() {
		toVar := mergeStmt.Stmt.(*stmt.AssignBinaryStmt).Lvalue.(expr.Var)
		mergeObj := s.HM.GetMergedObj(mergeStmt.Stmt)
		toP := s.CSM.GetCSVar(mergeStmt.Ctx, toVar)
		for _, csObj := range diff.Set0.ToSlice() {
			if strings.Contains(csObj.String(), "taint") ||
				strings.Contains(csObj.String(), "decode") ||
				strings.Contains(csObj.String(), "replace") ||
				strings.Contains(csObj.String(), "uncompress") {
				//mergeObj.Tags = "taint"
				//heapCtx := solver.Selector.SelectHeapContext(nil, mergeObj)
				o := csObj.Obj.(*heap.ConstObj)
				o.Count++
				s.AddVarPointsTo2(toP, csObj)
				//solver.AddVarPointsTo(mergeStmt.Ctx, toVar, heapCtx, mergeObj)
			} else {
				heapCtx := s.Selector.SelectHeapContext(nil, mergeObj)
				m := s.CSM.GetCSObj(heapCtx, mergeObj)
				s.AddVarPointsTo2(toP, m)
			}
		}
	}
}

// a.func1() a.func2() a's diff = [a1,a2],will create new call edge
func (s *Solver) processCall(cvar *context.CSVarPointer, diff *context.PointsToSet) {
	//todo
	//1. find all relevant invoke called by cvar
	//2. build new call edge as WList's entry
	//3. change `This` var by the new diff recObj
	ctx := cvar.Context
	var callee *lang.PHPMethod
	for _, callSite := range cvar.GetInvokeStmts() {

		for _, csObj := range diff.Set0.ToSlice() {
			invokeExpr := callSite.Rvalue.(*expr.InvokeInstanceExp)
			newObj, ok := csObj.Obj.(*heap.NewObj)
			if !ok {
				continue
			}
			if class, ok := newObj.GetObjType().(*lang.ClassType); ok {
				callee = s.world.Hierarchy.LookupMethod(class.Class, invokeExpr.MethodRef.Name)
			}
			//todo
			if callee == nil {
				if s.world.Debug {
					if invokeExpr.MethodRef.Name != "init" && invokeExpr.MethodRef.Name != "__destruct" {
						common.Log().Infof("Callee is NULL:%v", invokeExpr.MethodRef.Name)
					}
				}
				continue
			}
			csCallSite := s.CSM.GetCSCallSite(ctx, callSite)
			calleeCtx := s.Selector.SelectContext1(csCallSite, csObj, callee)

			csCallee := s.CSM.GetCSMethod(calleeCtx, callee)

			callEdge := context.CallEdge{
				CallSite: csCallSite,
				Callee:   csCallee,
			}
			s.AddCallEdge(callEdge)

			//pass this
			ir := s.world.GetIR(callee)
			this := s.CSM.GetCSVar(calleeCtx, ir.GetThis())
			s.AddVarPointsTo2(this, csObj)
		}
	}
}

// expand Method entry, main or other Method in Ir
func (s *Solver) processCallEdge(callEdge context.CallEdge) {
	csCallee := callEdge.Callee
	if csCallee.Method.Source == nil {
		return
	}
	s.AddCSMethod(csCallee)

	callerCtx := callEdge.CallSite.Ctx
	//CallSite := callEdge.CallSite.CallSite.GetRValue().(*InvokeFunctionExp)
	callSite := callEdge.CallSite.CallSite
	invokeExpr := callSite.GetInvokeExpr()
	calleeCtx := csCallee.Ctx
	callee := csCallee.Method
	//build arg point edge
	ir := s.world.IRS[callee.GetSignature()]
	for i, fromArg := range invokeExpr.GetArgs() {
		toParam := ir.GetParamAt(i)
		from := s.CSM.GetCSVar(callerCtx, fromArg)
		to := s.CSM.GetCSVar(calleeCtx, toParam)
		s.AddPFGEdge(from, to, context.PARAMETER_PASSING)
	}
	//build return Value point edge
	lvar := callSite.Lvalue
	csLvar := s.CSM.GetCSVar(callerCtx, lvar)
	for _, v := range ir.GetReturns() {
		csFromVar := s.CSM.GetCSVar(calleeCtx, v)
		s.AddPFGEdge(csFromVar, csLvar, context.RETURN)
	}
}

// when Source has new point to, then call AddPFGEdge
// AddPFGEdge will be called in body Stmts visitor, Params and return Value point to
// a= [o1,o2,o3]
// b = [o4,o5]
// a = b (Source)
// addPointsTo: a ->[o4,o5]
func (s *Solver) AddPFGEdge(source context.Pointer, target context.Pointer, kind context.EdgeKind) {
	//todo
	//PFG update Source Pointer and Target Pointer into graph by PFG.addEdge(edge)
	//add Pointer entry to worklist by addPointsTo(Pointer Pointer, PointsToSet pts)
	//
	edge := context.PointerFlowEdge{
		Source: source,
		Target: target,
		Kind:   kind,
	}
	s.PFG.AddEdge(edge)
	var sourceSet *context.PointsToSet
	if source.GetPointsToSet() == nil {
		sourceSet = &context.PointsToSet{Set0: mapset.NewSet[*context.CSObj]()}
	} else {
		sourceSet = source.GetPointsToSet()
	}
	if !sourceSet.IsEmpty() {
		s.AddVarPointsTo1(target, sourceSet)
	}

}
func (s *Solver) AddCallEdge(edge context.CallEdge) {
	s.CallEdges = append(s.CallEdges, edge)
	s.WList.CallEdges.Put(edge)
}

// when has a new Method call, Process point to in This Method's body
func (s *Solver) AddCSMethod(csMethod *context.CSMethod) {
	method := csMethod.Method
	ir := s.world.GetIR(method)
	if ir == nil {
		return
	}
	s.OpcodeTransform(ir)
	s.AddStmts(csMethod, ir.GetStmts())
}
func MockInstanceInvoke(m *lang.PHPMethod, base expr.Var, args []expr.Var) stmt.InvokeStmt {
	invokeExpr := &expr.InvokeInstanceExp{
		MethodRef: lang.MethodRef{
			DeclaringClass: m.DeclaringClass,
			Name:           m.Name,
			IsDynamic:      false,
			ReturnType:     nil,
			Signature:      m.GetSignature(),
		},
		Args: args,
		Base: base,
	}
	return stmt.InvokeStmt{
		Lvalue:     expr.Var{},
		Rvalue:     invokeExpr,
		Container:  nil,
		LineNumber: 0,
	}
}

func MockStaticInvoke(m *lang.PHPMethod, args []expr.Var) stmt.InvokeStmt {
	invokeExpr := &expr.InvokeStaticExp{
		FunRef: lang.MethodRef{
			DeclaringClass: m.DeclaringClass,
			Name:           m.Name,
			IsDynamic:      false,
			ReturnType:     nil,
			Signature:      m.GetSignature(),
		},
		Args:   args,
		FunVar: expr.Var{},
	}
	return stmt.InvokeStmt{
		Lvalue:     expr.Var{},
		Rvalue:     invokeExpr,
		Container:  nil,
		LineNumber: 0,
	}
}

// OpcodeTransform transfer magic syntax to normal ir opcode
func (s *Solver) OpcodeTransform(ir ir.IR) {
	for _, s0 := range ir.GetStmts() {
		switch stmt := s0.(type) {
		case *stmt.InvokeStmt:
			if invoke, ok := stmt.Rvalue.(*expr.InvokeStaticExp); ok {
				//callback function
				if invoke.FunRef.Name != "" && strings.Contains("call_user_func,array_map,ob_start", invoke.FunRef.Name) {
					if len(invoke.Args) == 0 {
						continue
					}
					callback := invoke.Args[0]
					parms := invoke.Args[1:]
					switch t := callback.Type.(type) {
					case *lang.ClassType:
						m := s.world.Hierarchy.LookupMethod(t.Class, "__invoke")
						//if m == nil {
						//	return
						//}
						instanceInvoke := MockInstanceInvoke(m, callback, parms)
						ir.AddStmt(&instanceInvoke)
					case *lang.ScalarStringType:
						callback := s.world.S.VM.GetConst(callback)
						cbName := callback.(string)
						cbName = common.TrimString(cbName)
						method := s.world.Hierarchy.LookupMethod(s.world.MainMethod.DeclaringClass, cbName)
						if method == nil && s.world.Debug {
							common.Log().Infof("No method:%v", cbName)
						}
						if invoke.FunRef.Name == "ob_start" {
							arg := expr.Var{
								Name: "global_buffer",
								Type: nil,
							}
							parms = append(parms, arg)
							stmt0 := MockStaticInvoke(method, parms)
							ir.AddStmt(&stmt0)
						} else {
							stmt0 := MockStaticInvoke(method, parms)
							ir.AddStmt(&stmt0)
						}
					default:
						//logrus.Debugf("No method:%v", callback.Type.String())
					}
				}
			}
		}
	}
}

func (s *Solver) AddStmts(csMethod *context.CSMethod, stmts []stmt.Stmt) {
	s.Process(csMethod, stmts)
}

func (s *Solver) AddVarPointsTo(ctx context.Context, var0 expr.Var, heapCtx context.Context, obj heap.Obj) {
	s.AddVarPointsTo2(s.CSM.GetCSVar(ctx, var0), s.CSM.GetCSObj(heapCtx, obj))
}
func (s *Solver) AddVarPointsTo1(pointer context.Pointer, pts *context.PointsToSet) {
	s.WList.AddEntry(pointer, pts)
}
func (s *Solver) AddVarPointsTo2(pointer context.Pointer, obj *context.CSObj) {
	pts := s.MakePointsToSet()
	pts.AddObj(obj)
	s.AddVarPointsTo1(pointer, &pts)
}
func (s *Solver) MakePointsToSet() context.PointsToSet {
	return context.PointsToSet{
		Set0: mapset.NewSet[*context.CSObj](),
	}
}
func (s *Solver) propagate(pointer context.Pointer, pts *context.PointsToSet) *context.PointsToSet {
	var diff *context.PointsToSet
	diff = pointer.Diff(pts)
	if len(diff.GetObjects()) > 0 {
		for _, v := range pointer.GetOutEdges() {
			edge := v.GetTarget().(context.Pointer)
			s.AddVarPointsTo1(edge, diff)
		}
	}
	return diff
}

func (s *Solver) Process(csMethod *context.CSMethod, stmts []stmt.Stmt) {
	visitor := CreateStmtVisitor(csMethod, s)
	for _, stmt := range stmts {
		stmt.Accept(visitor)
	}
}
