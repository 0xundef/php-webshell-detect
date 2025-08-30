package core

import (
	"github.com/0xundef/php-webshell-detect/internal/core/context"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
)

type DefaultStmtVisitor struct {
	stmt.StmtVisitor
	CsMethod *context.CSMethod
	Context  context.Context
	Solver   *Solver
}

func CreateStmtVisitor(method *context.CSMethod, solver *Solver) stmt.StmtVisitor {
	return &DefaultStmtVisitor{
		CsMethod: method,
		Context:  method.Ctx,
		Solver:   solver,
	}
}

func (v *DefaultStmtVisitor) VisitAssignAssignBinaryStmt(stmt stmt.AssignBinaryStmt) {
	var op1, op2 expr.Var
	to := stmt.Lvalue.(expr.Var)
	toPtr := v.Solver.CSM.GetCSVar(v.Context, to)

	op1 = stmt.Rvalue.(expr.BinaryExpr).GetLOperand()
	from1 := v.Solver.CSM.GetCSVar(v.Context, op1)
	from1.AddMerge(stmt, v.Context)

	op2 = stmt.Rvalue.(expr.BinaryExpr).GetROperand()
	from2 := v.Solver.CSM.GetCSVar(v.Context, op2)
	from2.AddMerge(stmt, v.Context)

	toPtr.MergeCount = from1.MergeCount + from2.MergeCount + 1
}

// directly joint to Pointer->Set0 by AddVarPointsTo
func (v *DefaultStmtVisitor) VisitAssignLiteralStmt(stmt stmt.AssignLiteralStmt) {
	rv := stmt.GetRValue()
	if lobj, ok := rv.(expr.Literal); ok {
		obj := v.Solver.HM.GetConstObj(lobj)
		heapCtx := v.Solver.Selector.SelectHeapContext(v.CsMethod, obj)
		var0 := stmt.GetLValue().(expr.Var)
		v.Solver.AddVarPointsTo(v.Context, var0, heapCtx, obj)
	}
}

// new point,AddPFGEdge
func (v *DefaultStmtVisitor) VisitAssignVarStmt(stmt stmt.AssignVarStmt) {
	lv := stmt.GetLValue().(expr.Var)
	rv := stmt.GetRValue().(expr.Var)
	from := v.Solver.CSM.GetCSVar(v.Context, rv)
	to := v.Solver.CSM.GetCSVar(v.Context, lv)
	to.SetMergeCount(from.GetMergeCount())
	v.Solver.AddPFGEdge(from, to, context.LOCAL)
}

// will generate call edges, as entry of work list
func (v *DefaultStmtVisitor) VisitInvokeStaticStmt(callSite stmt.InvokeStmt) {
	if callSite.Container == nil {
		callSite.Container = v.CsMethod.Method
	}
	if callSite.Container.Name == callSite.Rvalue.GetFunName() {
		return
	}
	if invokeExpr, ok := callSite.Rvalue.(*expr.InvokeStaticExp); ok {

		callee := v.Solver.world.Hierarchy.ResolveMethod(invokeExpr.FunRef)
		//function not found, build-in function
		//todo
		if callee == nil {
			var args []string
			for _, arg := range invokeExpr.Args {
				args = append(args, arg.Name)
			}
			var funName string
			if invokeExpr.FunVar.Name != "" {
				dynamicFun := v.Solver.CSM.GetCSVar(v.Context, invokeExpr.FunVar)
				dynamicFun.IsFunName = true
				funName = "var_fun"
				csLvar := v.Solver.CSM.GetCSVar(v.Context, callSite.Lvalue)
				csLvar.SetMergeCount(dynamicFun.GetMergeCount())
			} else {
				funName = invokeExpr.FunRef.Name
			}
			callee = &lang.PHPMethod{
				DeclaringClass: nil,
				Name:           funName,
				ParamTypes:     nil,
				ReturnType:     nil,
				ParamNames:     args,
				Source:         nil,
			}
		}
		if invokeExpr.FunVar.Name != "" {
			dynamicFun := v.Solver.CSM.GetCSVar(v.Context, invokeExpr.FunVar)
			dynamicFun.IsFunName = true
			csLvar := v.Solver.CSM.GetCSVar(v.Context, callSite.Lvalue)
			csLvar.SetMergeCount(dynamicFun.GetMergeCount())
		}
		csCallSite := v.Solver.CSM.GetCSCallSite(v.Context, callSite)
		calleeCtx := v.Solver.Selector.SelectContext0(csCallSite, callee)
		csCallee := v.Solver.CSM.GetCSMethod(calleeCtx, callee)
		callEdge := context.CallEdge{
			CallSite: csCallSite,
			Callee:   csCallee,
		}
		v.Solver.AddCallEdge(callEdge)
	}
	if invokeExpr, ok := callSite.Rvalue.(*expr.InvokeInstanceExp); ok {
		p := v.Solver.CSM.GetCSVar(v.Context, invokeExpr.Base)
		p.AddInvokeStmt(callSite)
	}
}

func (v *DefaultStmtVisitor) VisitNewStmt(newSite stmt.NewStmt) {
	newObj := v.Solver.HM.GetNewObj(newSite)
	heapCtx := v.Solver.Selector.SelectHeapContext(v.CsMethod, newObj)
	lvar := newSite.Lvalue
	v.Solver.AddVarPointsTo(v.Context, lvar, heapCtx, newObj)
}

func (v *DefaultStmtVisitor) VisitLoadFieldStmt(load stmt.LoadFieldStmt) {
	if !load.Rvalue.FieldRef.IsStatic() {
		p := v.Solver.CSM.GetCSVar(v.Context, load.Rvalue.Base)
		p.AddLoadField(load)
	}
}

func (v *DefaultStmtVisitor) VisitStoreFieldStmt(store stmt.StoreFieldStmt) {
	if !store.Lvalue.FieldRef.IsStatic() {
		p := v.Solver.CSM.GetCSVar(v.Context, store.Lvalue.Base)
		p.AddStoreField(store)
	}
}

func (v *DefaultStmtVisitor) VisitStoreArrayStmt(stmt stmt.StoreArrayStmt) {
	base := v.Solver.CSM.GetCSVar(v.Context, stmt.Lvalue.Base)
	base.AddStoreArray(stmt)
}

// VisitLoadArrayStmt $a = $b[0]    stmt.Rvalue.Base--->$b
func (v *DefaultStmtVisitor) VisitLoadArrayStmt(stmt stmt.LoadArrayStmt) {
	base := v.Solver.CSM.GetCSVar(v.Context, stmt.Rvalue.Base)
	base.AddLoadArray(stmt, v.Context)
	if len(base.PointsToSet.GetObjects()) > 0 {
		tP := v.Solver.CSM.GetCSVar(v.Context, stmt.Lvalue)
		for _, csObj := range base.PointsToSet.GetObjects() {
			sP := v.Solver.CSM.GetArrayPointer(csObj)
			v.Solver.AddPFGEdge(sP, tP, 1)
		}
	}
}
