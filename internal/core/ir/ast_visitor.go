package ir

import (
	"github.com/0xundef/php-webshell-detect/internal/core/common"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/expr"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/lang"
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"github.com/0xundef/php-webshell-detect/frontend/php/pkg/ast"
	"strconv"
	"strings"
)

// AstVisitor the whole stage of the parser frontend, the output result is IR, it is not necessary to implement all the Vertex visit method
type AstVisitor struct {
	Ctx        VisitContext
	Ir         IR
	Hierarchy  *lang.PHPClassHierarchy
	BuildClass bool
	VM         *VarManager
	SeqIndex   int
	Debug      bool
}

// VisitContext this context will return the child IR element at current parent Vertex's visit method
type VisitContext struct {
	ret    interface{}
	parent ast.Vertex
}

// visitor will drive the child Vertex and return IR element by VisitContext
func (visitor *AstVisitor) drive(n ast.Vertex) {
	if n != nil {
		n.Accept(visitor)
	}
}
func (visitor *AstVisitor) drive1(n ast.Vertex, current ast.Vertex, parent ast.Vertex) {
	if n != nil {
		visitor.Ctx.ret = nil
		visitor.Ctx.parent = current
		n.Accept(visitor)
		visitor.Ctx.parent = parent
	}
}

func (visitor *AstVisitor) driveList(list []ast.Vertex) {
	for _, nn := range list {
		if visitor.BuildClass {
			if _, ok := nn.(*ast.StmtClass); ok {
				visitor.drive1(nn, nil, nil)
			}
		} else {
			visitor.drive1(nn, nil, nil)
		}
	}
}

func (visitor *AstVisitor) Root(n *ast.Root) {
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) Nullable(n *ast.Nullable) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) Parameter(n *ast.Parameter) {
	visitor.drive(n.Type)
	visitor.drive(n.Var)
	//todo fix
	//visitor.drive(n.DefaultValue)
}

func (visitor *AstVisitor) Identifier(n *ast.Identifier) {
	visitor.Ctx.ret = string(n.Value)
}

func (visitor *AstVisitor) Argument(n *ast.Argument) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) StmtBreak(n *ast.StmtBreak) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) StmtCase(n *ast.StmtCase) {
	visitor.drive(n.Cond)
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtCatch(n *ast.StmtCatch) {
	visitor.drive(n.Var)
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtClass(n *ast.StmtClass) {
	if !visitor.BuildClass {
		return
	}
	className := ConvertString(n.Name)
	class := lang.CreatePHPClass(className, "")
	var initNodes []ast.Vertex
	for _, stmt := range n.Stmts {
		if stmtMethod, ok := stmt.(*ast.StmtClassMethod); ok {
			bodyStmts := stmtMethod.Stmt.(*ast.StmtStmtList)
			source := ast.Root{
				Position: n.Position,
				Stmts:    bodyStmts.Stmts,
				EndTkn:   nil,
			}
			var params []string
			for _, param := range stmtMethod.Params {
				visitor.drive1(param, n, visitor.Ctx.parent)
				v := visitor.Ctx.ret.(expr.Var)
				params = append(params, v.Name)
			}
			name := ConvertString(stmtMethod.Name)
			m := &lang.PHPMethod{
				Name:           name,
				DeclaringClass: class,
				ParamTypes:     nil,
				ReturnType:     nil,
				ParamNames:     params,
				Source:         &source,
			}
			class.AddMethod(m)
		}
		if properties, ok := stmt.(*ast.StmtPropertyList); ok {
			//todo modifiers, Type
			for _, node := range properties.Props {
				var name string
				if prop, ok := node.(*ast.StmtProperty); ok {
					if varExpr, ok := prop.Var.(*ast.ExprVariable); ok {
						name = ConvertString(varExpr.Name)
					}
				}
				if name != "" {
					f := &lang.PHPClassField{
						DeclaringClass: class,
						Name:           name,
						Type0:          nil,
					}
					class.AddField(f)
				}
				initNodes = append(initNodes, node)
			}
		}
	}
	//mock init constructor
	initSource := ast.Root{
		Position: n.Position,
		Stmts:    initNodes,
		EndTkn:   nil,
	}
	init := &lang.PHPMethod{
		DeclaringClass: class,
		Name:           "init",
		ParamTypes:     nil,
		ReturnType:     nil,
		ParamNames:     nil,
		Source:         &initSource,
	}
	class.AddMethod(init)
	visitor.Hierarchy.AddClass(class)
}

func (visitor *AstVisitor) StmtClassConstList(n *ast.StmtClassConstList) {
	visitor.driveList(n.Modifiers)
}

func (visitor *AstVisitor) StmtClassMethod(n *ast.StmtClassMethod) {
	visitor.driveList(n.Modifiers)
	visitor.drive(n.Name)
	visitor.drive(n.ReturnType)
	visitor.drive(n.Stmt)
}

func (visitor *AstVisitor) StmtConstList(n *ast.StmtConstList) {
}

func (visitor *AstVisitor) StmtConstant(n *ast.StmtConstant) {
	visitor.drive(n.Name)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) StmtContinue(n *ast.StmtContinue) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) StmtDeclare(n *ast.StmtDeclare) {
	if stmt, ok := n.Stmt.(*ast.StmtStmtList); ok && n.ColonTkn != nil {
		visitor.driveList(stmt.Stmts)
	} else {
		visitor.drive(n.Stmt)
	}
}

func (visitor *AstVisitor) StmtDefault(n *ast.StmtDefault) {
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtDo(n *ast.StmtDo) {
	visitor.drive(n.Stmt)
	visitor.drive(n.Cond)
}

func (visitor *AstVisitor) StmtEcho(n *ast.StmtEcho) {
	for _, vertex := range n.Exprs {
		visitor.drive1(vertex, n, visitor.Ctx.parent)
		if rvalue, ok := visitor.Ctx.ret.(expr.Var); ok {
			lv := expr.Var{
				Name: "global_buffer",
				Type: nil,
			}
			e := &stmt.AssignVarStmt{
				Lvalue:     lv,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
	}
}

func (visitor *AstVisitor) StmtElse(n *ast.StmtElse) {
	if stmt, ok := n.Stmt.(*ast.StmtStmtList); ok && n.ColonTkn != nil {
		visitor.driveList(stmt.Stmts)
	} else {
		visitor.drive(n.Stmt)
	}
}

func (visitor *AstVisitor) StmtElseIf(n *ast.StmtElseIf) {
	visitor.drive(n.Cond)
	if stmt, ok := n.Stmt.(*ast.StmtStmtList); ok && n.ColonTkn != nil {
		visitor.driveList(stmt.Stmts)
	} else {
		visitor.drive(n.Stmt)
	}
}

func (visitor *AstVisitor) StmtExpression(n *ast.StmtExpression) {
	visitor.drive1(n.Expr, n, visitor.Ctx.parent)
}

func (visitor *AstVisitor) StmtFinally(n *ast.StmtFinally) {
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtFor(n *ast.StmtFor) {
	if stmt, ok := n.Stmt.(*ast.StmtStmtList); ok && n.ColonTkn != nil {
		visitor.driveList(stmt.Stmts)
	} else {
		visitor.drive(n.Stmt)
	}
}

func (visitor *AstVisitor) StmtForeach(n *ast.StmtForeach) {
	visitor.drive(n.Expr)
	visitor.drive(n.Key)
	visitor.drive(n.Var)
	if stmt, ok := n.Stmt.(*ast.StmtStmtList); ok && n.ColonTkn != nil {
		visitor.driveList(stmt.Stmts)
	} else {
		visitor.drive(n.Stmt)
	}
}

// function definition
func (visitor *AstVisitor) StmtFunction(n *ast.StmtFunction) {
	//todo return type
	var declaringClass *lang.PHPClass
	var name string
	var params []string
	visitor.drive1(n.Name, n, visitor.Ctx.parent)
	name = visitor.Ctx.ret.(string)
	for _, param := range n.Params {
		visitor.drive1(param, n, visitor.Ctx.parent)
		v := visitor.Ctx.ret.(expr.Var)
		params = append(params, v.Name)
	}

	source := ast.Root{
		Position: n.Position,
		Stmts:    n.Stmts,
		EndTkn:   nil,
	}
	declaringClass = visitor.Hierarchy.GetClass("EntryMain")
	m := &lang.PHPMethod{
		DeclaringClass: declaringClass,
		Name:           name,
		ParamTypes:     nil,
		ReturnType:     nil,
		ParamNames:     params,
		Source:         &source,
	}
	declaringClass.AddMethod(m)
}

func (visitor *AstVisitor) StmtGlobal(n *ast.StmtGlobal) {
	for _, vertex := range n.Vars {
		visitor.drive1(vertex, n, visitor.Ctx.parent)
		switch v := visitor.Ctx.ret.(type) {
		case expr.ArrayAccess:
		case expr.Var:
			v.Global = true
			visitor.VM.AddGlobalVar(v.Name, v)
		default:
		}
	}
}

func (visitor *AstVisitor) StmtGoto(n *ast.StmtGoto) {
	visitor.drive(n.Label)
}

func (visitor *AstVisitor) StmtHaltCompiler(n *ast.StmtHaltCompiler) {
}

func (visitor *AstVisitor) StmtIf(n *ast.StmtIf) {
	visitor.drive(n.Cond)
	if stmt, ok := n.Stmt.(*ast.StmtStmtList); ok && n.ColonTkn != nil {
		visitor.driveList(stmt.Stmts)
	} else {
		visitor.drive(n.Stmt)
	}
	visitor.driveList(n.ElseIf)
	visitor.drive(n.Else)
}

func (visitor *AstVisitor) StmtInlineHtml(n *ast.StmtInlineHtml) {
	html := string(n.Value)
	inputs := common.ExtractInputs(html)
	for _, input := range inputs {
		lV := expr.Var{
			Name:   "$" + input,
			Type:   &lang.ScalarStringType{},
			Global: false,
		}
		rV := expr.ScalarString{
			Value: "taint",
		}
		st := &stmt.AssignLiteralStmt{
			Lvalue:     lV,
			Rvalue:     &rV,
			LineNumber: n.Position.StartLine,
			Seq:        0,
		}
		visitor.AddStmt(st)
	}
}

func (visitor *AstVisitor) StmtInterface(n *ast.StmtInterface) {
	visitor.drive(n.Name)
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtLabel(n *ast.StmtLabel) {
	visitor.drive(n.Name)
}

func (visitor *AstVisitor) StmtNamespace(n *ast.StmtNamespace) {
	visitor.drive(n.Name)
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtNop(n *ast.StmtNop) {
}

func (visitor *AstVisitor) StmtProperty(n *ast.StmtProperty) {
	visitor.drive1(n.Var, n, visitor.Ctx.parent)
	fV := visitor.Ctx.ret.(expr.Var)
	class := visitor.Ir.GetContainerClass()
	classField := visitor.Hierarchy.LookupField(class, fV.Name)
	field := expr.FieldAccess{
		Base: visitor.Ir.GetThis(),
		FieldRef: lang.FieldRef{
			DeclaringClass: *class,
			Name:           fV.Name,
			Field:          *classField,
			Type0:          nil,
			IsStatic0:      false,
		},
	}
	visitor.drive1(n.Expr, n, visitor.Ctx.parent)
	rvalue := visitor.Ctx.ret

	switch v := rvalue.(type) {
	case expr.Var:
		s2 := &stmt.StoreFieldStmt{
			Lvalue:     field,
			Rvalue:     v,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(s2)
	case expr.NewInstanceExpr:
		tmp := MakeTmpVar(n, nil)
		s1 := &stmt.NewStmt{
			Lvalue:     tmp,
			Rvalue:     &v,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(s1)

		s2 := &stmt.StoreFieldStmt{
			Lvalue:     field,
			Rvalue:     tmp,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(s2)
	}
}

func (visitor *AstVisitor) StmtPropertyList(n *ast.StmtPropertyList) {
	visitor.driveList(n.Modifiers)
	visitor.drive(n.Type)
}

func (visitor *AstVisitor) StmtReturn(n *ast.StmtReturn) {
	visitor.drive1(n.Expr, n, visitor.Ctx.parent)
	switch ret := visitor.Ctx.ret.(type) {
	case expr.Var:
		stmt := &stmt.ReturnStmt{
			Value:      ret,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(stmt)
		visitor.Ir.AddReturn(ret)
		break
	case expr.FieldAccess:
		tmp := expr.Var{
			Name: "%ret" + strconv.Itoa(n.Position.StartPos) + "_" + strconv.Itoa(n.Position.EndPos),
			Type: nil,
		}
		s1 := &stmt.LoadFieldStmt{
			Lvalue:     tmp,
			Rvalue:     ret,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(s1)

		stmt := &stmt.ReturnStmt{
			Value:      tmp,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(stmt)
		visitor.Ir.AddReturn(tmp)
		break
	default:
		if visitor.Debug {
			common.Log().Errorf("StmtReturn:%v", ret)
		}
	}

}

func (visitor *AstVisitor) StmtStatic(n *ast.StmtStatic) {
}

func (visitor *AstVisitor) StmtStaticVar(n *ast.StmtStaticVar) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) StmtStmtList(n *ast.StmtStmtList) {
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtSwitch(n *ast.StmtSwitch) {
	visitor.drive(n.Cond)
	visitor.driveList(n.Cases)
}

func (visitor *AstVisitor) StmtThrow(n *ast.StmtThrow) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) StmtTrait(n *ast.StmtTrait) {
	visitor.drive(n.Name)
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) StmtTraitUse(n *ast.StmtTraitUse) {
	visitor.driveList(n.Adaptations)
}

func (visitor *AstVisitor) StmtTraitUseAlias(n *ast.StmtTraitUseAlias) {
	visitor.drive(n.Trait)
	visitor.drive(n.Method)
	visitor.drive(n.Modifier)
	visitor.drive(n.Alias)
}

func (visitor *AstVisitor) StmtTraitUsePrecedence(n *ast.StmtTraitUsePrecedence) {
	visitor.drive(n.Trait)
	visitor.drive(n.Method)
}

func (visitor *AstVisitor) StmtTry(n *ast.StmtTry) {
	visitor.driveList(n.Stmts)
	visitor.driveList(n.Catches)
	visitor.drive(n.Finally)
}

func (visitor *AstVisitor) StmtUnset(n *ast.StmtUnset) {
}

func (visitor *AstVisitor) StmtUse(n *ast.StmtUseList) {
	visitor.drive(n.Type)
}

func (visitor *AstVisitor) StmtGroupUse(n *ast.StmtGroupUseList) {
	visitor.drive(n.Type)
	visitor.drive(n.Prefix)
}

func (visitor *AstVisitor) StmtUseDeclaration(n *ast.StmtUse) {
	visitor.drive(n.Type)
	visitor.drive(n.Use)
	visitor.drive(n.Alias)
}

func (visitor *AstVisitor) StmtWhile(n *ast.StmtWhile) {
	visitor.drive(n.Cond)
	if stmt, ok := n.Stmt.(*ast.StmtStmtList); ok && n.ColonTkn != nil {
		visitor.driveList(stmt.Stmts)
	} else {
		visitor.drive(n.Stmt)
	}
}

func (visitor *AstVisitor) ExprArray(n *ast.ExprArray) {
	exprNew := expr.NewInstanceExpr{
		Type0: lang.ClassType{
			Name:  "Array",
			Class: nil,
		},
	}
	visitor.Ctx.ret = exprNew
}

func (visitor *AstVisitor) ExprArrayDimFetch(n *ast.ExprArrayDimFetch) {
	var arrayName string
	visitor.drive(n.Var)
	if v, ok := visitor.Ctx.ret.(expr.Var); ok {
		arrayName = v.Name
	}
	var pos int
	switch v := visitor.Ctx.parent.(type) {
	case *ast.ExprAssign:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignCoalesce:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignConcat:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignDiv:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignBitwiseOr:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignBitwiseAnd:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignBitwiseXor:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignMinus:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignMod:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignMul:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignPlus:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignPow:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignReference:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignShiftRight:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignShiftLeft:
		pos = v.Var.GetPosition().StartPos
	default:
		pos = -1
	}
	var global bool
	if strings.Contains("$_GET,$_POST,$_COOKIE,$_REQUEST", arrayName) {
		global = true
	}
	// lvalue
	if pos == n.Position.StartPos {
		visitor.Ctx.ret = expr.ArrayAccess{
			Base: expr.Var{
				Name:   arrayName,
				Type:   nil,
				Global: global,
			},
			Index: 0,
		}
	} else {
		a := expr.ArrayAccess{
			Base: expr.Var{
				Name:   arrayName,
				Type:   nil,
				Global: global,
			},
			Index: 0,
		}

		ret := MakeTmpVar(n, a.GetType())
		st := &stmt.LoadArrayStmt{
			Lvalue:     ret,
			Rvalue:     a,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(st)
		visitor.Ctx.ret = ret
	}
}

func (visitor *AstVisitor) ExprArrayItem(n *ast.ExprArrayItem) {
	visitor.drive(n.Key)
	visitor.drive(n.Val)
}

func (visitor *AstVisitor) ExprArrowFunction(n *ast.ExprArrowFunction) {
	visitor.drive(n.ReturnType)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprBitwiseNot(n *ast.ExprBitwiseNot) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprBooleanNot(n *ast.ExprBooleanNot) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprBrackets(n *ast.ExprBrackets) {
	visitor.drive1(n.Expr, n, visitor.Ctx.parent)
}

func (visitor *AstVisitor) ExprClassConstFetch(n *ast.ExprClassConstFetch) {
	visitor.drive(n.Class)
	visitor.drive(n.Const)
}

func (visitor *AstVisitor) ExprClone(n *ast.ExprClone) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprClosure(n *ast.ExprClosure) {
	visitor.drive(n.ReturnType)
	visitor.driveList(n.Stmts)
}

func (visitor *AstVisitor) ExprClosureUse(n *ast.ExprClosureUse) {
	visitor.drive(n.Var)
}

func (visitor *AstVisitor) ExprConstFetch(n *ast.ExprConstFetch) {
	visitor.drive1(n.Const, n, visitor.Ctx.parent)
	value := visitor.Ctx.ret
	rv := &expr.ScalarString{
		Value: value.(string),
	}
	tmp := MakeTmpVar(n, nil)
	stmt := &stmt.AssignLiteralStmt{
		Lvalue:     tmp,
		Rvalue:     rv,
		LineNumber: n.Position.StartLine,
	}
	visitor.AddStmt(stmt)
	visitor.Ctx.ret = tmp
}

func (visitor *AstVisitor) ExprEmpty(n *ast.ExprEmpty) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprErrorSuppress(n *ast.ExprErrorSuppress) {
	visitor.drive1(n.Expr, n, visitor.Ctx.parent)
}

func (visitor *AstVisitor) ExprEval(n *ast.ExprEval) {
	visitor.drive1(n.Expr, n, visitor.Ctx.parent)
	var args []expr.Var
	switch code := visitor.Ctx.ret.(type) {
	case ast.ScalarEncapsed:
	case ast.ExprFunctionCall:
	case expr.Var:
		declaringClass := visitor.Hierarchy.GetClass("EntryMain") //todo remove
		sig := declaringClass.Name + "::eval_fake"
		args = append(args, code)
		call := &expr.InvokeStaticExp{
			FunRef: lang.MethodRef{
				DeclaringClass: declaringClass,
				Name:           "eval_fake",
				IsDynamic:      false,
				ReturnType:     nil,
				Signature:      sig,
			},
			Args:   args,
			FunVar: expr.Var{},
		}
		lvar := MakeTmpVar(n, &lang.UnKnowType{})
		stmt := &stmt.InvokeStmt{
			Lvalue:     lvar,
			Rvalue:     call,
			Container:  visitor.Ir.GetContainerMethod(),
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(stmt)
		visitor.Ctx.ret = lvar
	default:
		if visitor.Debug {
			common.Log().Errorf("ExprEval:%v", code)
		}
	}
}

func (visitor *AstVisitor) ExprExit(n *ast.ExprExit) {
	visitor.drive(n.Expr)
}

// pure function call
func (visitor *AstVisitor) ExprFunctionCall(n *ast.ExprFunctionCall) {
	visitor.drive1(n.Function, n, visitor.Ctx.parent)
	fun := visitor.Ctx.ret
	args := []expr.Var{}
	var rv expr.Invoke
	for i := 0; i < len(n.Args); i++ {
		visitor.drive1(n.Args[i], n, visitor.Ctx.parent)
		arg := visitor.MyArgTransfer(visitor.Ctx.ret, n.Args[i])
		args = append(args, arg)
	}
	declaringClass := visitor.Hierarchy.GetClass("EntryMain") //todo remove
	switch v := fun.(type) {
	case string:
		sig := declaringClass.Name + "::" + v
		rv = &expr.InvokeStaticExp{
			FunRef: lang.MethodRef{
				DeclaringClass: declaringClass,
				Name:           v,
				IsDynamic:      false,
				ReturnType:     nil,
				Signature:      sig,
			},
			Args: args,
		}
	case expr.Var:
		sig := declaringClass.Name + "::var_fun"
		rv = &expr.InvokeStaticExp{
			FunVar: v,
			FunRef: lang.MethodRef{
				DeclaringClass: declaringClass,
				Name:           "var_fun",
				IsDynamic:      false,
				ReturnType:     nil,
				Signature:      sig,
			},
			Args: args,
		}
	}

	lv := MakeTmpVar(n, rv.GetType())
	s := stmt.InvokeStmt{
		Lvalue:     lv,
		Rvalue:     rv,
		LineNumber: n.Position.StartLine,
	}
	visitor.AddStmt(&s)
	visitor.Ctx.ret = lv
}

func (visitor *AstVisitor) ExprInclude(n *ast.ExprInclude) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprIncludeOnce(n *ast.ExprIncludeOnce) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprInstanceOf(n *ast.ExprInstanceOf) {
	visitor.drive(n.Expr)
	visitor.drive(n.Class)
}

func (visitor *AstVisitor) ExprIsset(n *ast.ExprIsset) {
}

func (visitor *AstVisitor) ExprList(n *ast.ExprList) {
}
func (visitor *AstVisitor) MyArgTransfer(param interface{}, v ast.Vertex) expr.Var {
	var ret expr.Var
	switch arg := param.(type) {
	case expr.Var:
		ret = arg
		break
	case expr.NewInstanceExpr:
		lv := MakeTmpVar(v, &arg.Type0)
		newStmt := &stmt.NewStmt{
			Lvalue:     lv,
			Rvalue:     &arg,
			LineNumber: v.GetPosition().StartLine,
		}
		visitor.AddStmt(newStmt)
		ret = lv
		break
	case expr.ArrayAccess:
		tmp := MakeTmpVar(v, arg.GetType())
		stmt := &stmt.LoadArrayStmt{
			Lvalue:     tmp,
			Rvalue:     arg,
			LineNumber: v.GetPosition().StartLine,
		}
		visitor.AddStmt(stmt)
		ret = tmp
	case expr.FieldAccess:
		tmp := MakeTmpVar(v, arg.GetType())
		stmt := &stmt.LoadFieldStmt{
			Lvalue:     tmp,
			Rvalue:     arg,
			LineNumber: v.GetPosition().StartLine,
		}
		visitor.AddStmt(stmt)
		ret = tmp
	default:
		if visitor.Debug {
			common.Log().Errorf("MyArgTransfer Arg:%v  not match", arg)
		}
	}
	return ret
}

func (visitor *AstVisitor) ExprMethodCall(n *ast.ExprMethodCall) {
	//build Base
	visitor.drive1(n.Var, n, visitor.Ctx.parent)
	base := visitor.Ctx.ret.(expr.Var)

	visitor.drive1(n.Method, n, visitor.Ctx.parent)
	fun := visitor.Ctx.ret
	args := []expr.Var{}
	var rv expr.Invoke
	for i := 0; i < len(n.Args); i++ {
		visitor.drive1(n.Args[i], n, visitor.Ctx.parent)
		arg := visitor.MyArgTransfer(visitor.Ctx.ret, n.Args[i])
		args = append(args, arg)
	}
	if mName, ok := fun.(string); ok {
		rv = &expr.InvokeInstanceExp{
			MethodRef: lang.MethodRef{
				DeclaringClass: nil,
				Name:           mName,
				IsDynamic:      false,
				ReturnType:     nil,
				Signature:      mName,
			},
			Args: args,
			Base: base,
		}
	}

	lv := MakeTmpVar(n, rv.GetType())
	s := stmt.InvokeStmt{
		Lvalue:     lv,
		Rvalue:     rv,
		LineNumber: n.Position.StartLine,
	}
	visitor.AddStmt(&s)
	visitor.Ctx.ret = lv
}

func (visitor *AstVisitor) ExprNew(n *ast.ExprNew) {
	visitor.drive1(n.Class, n, visitor.Ctx.parent)
	className := visitor.Ctx.ret.(string)
	class := visitor.Hierarchy.GetClass(className)
	exprNew := expr.NewInstanceExpr{
		Type0: lang.ClassType{
			Name:  className,
			Class: class,
		},
	}
	visitor.Ctx.ret = exprNew
}

func (visitor *AstVisitor) ExprPostDec(n *ast.ExprPostDec) {
	visitor.drive(n.Var)
}

func (visitor *AstVisitor) ExprPostInc(n *ast.ExprPostInc) {
	visitor.drive(n.Var)
}

func (visitor *AstVisitor) ExprPreDec(n *ast.ExprPreDec) {
	visitor.drive(n.Var)
}

func (visitor *AstVisitor) ExprPreInc(n *ast.ExprPreInc) {
	visitor.drive(n.Var)
}

func (visitor *AstVisitor) ExprPrint(n *ast.ExprPrint) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprPropertyFetch(n *ast.ExprPropertyFetch) {
	var base expr.Var
	var ref lang.FieldRef
	visitor.drive1(n.Var, n, visitor.Ctx.parent)
	base = visitor.Ctx.ret.(expr.Var)

	visitor.drive1(n.Prop, n, visitor.Ctx.parent)
	filedName := visitor.Ctx.ret
	ref = lang.FieldRef{
		DeclaringClass: lang.PHPClass{},
		Name:           "$" + filedName.(string),
		Field:          lang.PHPClassField{},
		Type0:          nil,
		IsStatic0:      false,
	}
	var pos int
	switch v := visitor.Ctx.parent.(type) {
	case *ast.ExprAssign:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignCoalesce:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignConcat:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignDiv:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignBitwiseOr:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignBitwiseAnd:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignBitwiseXor:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignMinus:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignMod:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignMul:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignPlus:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignPow:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignReference:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignShiftRight:
		pos = v.Var.GetPosition().StartPos
	case *ast.ExprAssignShiftLeft:
		pos = v.Var.GetPosition().StartPos
	default:
		pos = -1
	}
	if pos == n.Position.StartPos {
		visitor.Ctx.ret = expr.FieldAccess{
			Base:     base,
			FieldRef: ref,
		}
	} else {
		f := expr.FieldAccess{
			Base:     base,
			FieldRef: ref,
		}

		ret := MakeTmpVar(n, f.GetType())
		st := &stmt.LoadFieldStmt{
			Lvalue:     ret,
			Rvalue:     f,
			LineNumber: 0,
		}
		visitor.AddStmt(st)
		visitor.Ctx.ret = ret
	}
}

func (visitor *AstVisitor) ExprRequire(n *ast.ExprRequire) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprRequireOnce(n *ast.ExprRequireOnce) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprShellExec(n *ast.ExprShellExec) {
	args := []expr.Var{}
	var rv expr.Invoke
	for _, part := range n.Parts {
		visitor.drive1(part, n, visitor.Ctx.parent)
		arg := visitor.MyArgTransfer(visitor.Ctx.ret, part)
		args = append(args, arg)
	}

	declaringClass := visitor.Hierarchy.GetClass("EntryMain") //todo remove
	sig := declaringClass.Name + "::shell_exec_2quotes"
	rv = &expr.InvokeStaticExp{
		FunRef: lang.MethodRef{
			DeclaringClass: declaringClass,
			Name:           "shell_exec_2quotes",
			IsDynamic:      false,
			ReturnType:     nil,
			Signature:      sig,
		},
		Args: args,
	}
	lv := MakeTmpVar(n, rv.GetType())
	s := stmt.InvokeStmt{
		Lvalue:     lv,
		Rvalue:     rv,
		LineNumber: n.Position.StartLine,
	}
	visitor.AddStmt(&s)
}

func (visitor *AstVisitor) ExprStaticCall(n *ast.ExprStaticCall) {
	visitor.drive(n.Class)
	visitor.drive(n.Call)
}

func (visitor *AstVisitor) ExprStaticPropertyFetch(n *ast.ExprStaticPropertyFetch) {
	visitor.drive(n.Class)
	visitor.drive(n.Prop)
}

func (visitor *AstVisitor) ExprTernary(n *ast.ExprTernary) {
	visitor.drive1(n.Cond, n, visitor.Ctx.parent)

	visitor.drive1(n.IfTrue, n, visitor.Ctx.parent)
	visitor.drive1(n.IfFalse, n, visitor.Ctx.parent)
}

func (visitor *AstVisitor) ExprUnaryMinus(n *ast.ExprUnaryMinus) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprUnaryPlus(n *ast.ExprUnaryPlus) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprVariable(n *ast.ExprVariable) {
	var v expr.Var
	var global bool

	if id, ok := n.Name.(*ast.Identifier); ok {
		name := string(id.Value)
		if PHP_CONTEXT.IsGlobal(name) {
			global = true
		}
		gV := visitor.VM.GetGlobalVar(name)
		if gV != nil {
			v = gV.(expr.Var)
		} else {
			v = MakeVar(name, &lang.UnKnowType{}, global)
			if global {
				visitor.VM.AddGlobalVar(name, v)
			}
		}
		visitor.Ctx.ret = v
	} else {
		visitor.drive1(n.Name, n, visitor.Ctx.parent)
		if v0, ok := visitor.Ctx.ret.(expr.Var); ok {
			v = MakeVar("$"+v0.Name, &lang.UnKnowType{}, global)
			visitor.Ctx.ret = v
		}
	}
}

func (visitor *AstVisitor) ExprYield(n *ast.ExprYield) {
	visitor.drive(n.Key)
	visitor.drive(n.Val)
}

func (visitor *AstVisitor) ExprYieldFrom(n *ast.ExprYieldFrom) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssign(n *ast.ExprAssign) {
	visitor.drive1(n.Var, n, visitor.Ctx.parent)
	if v, ok := visitor.Ctx.ret.(expr.Var); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret

		if rvalue, ok := ret.(expr.BinaryExpr); ok {
			s := stmt.AssignBinaryStmt{
				DefinitionStmt: nil,
				Lvalue:         v,
				Rvalue:         rvalue,
				LineNumber:     n.Position.StartLine,
			}
			visitor.AddStmt(&s)
		}
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.AssignVarStmt{
				DefinitionStmt: nil,
				Lvalue:         v,
				Rvalue:         rvalue,
				LineNumber:     n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
		if rvalue, ok := ret.(expr.NewInstanceExpr); ok {
			e := &stmt.NewStmt{
				Lvalue:     v,
				Rvalue:     &rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
			//mock init
			invoke := expr.InvokeInstanceExp{
				MethodRef: lang.MethodRef{
					DeclaringClass: visitor.Ir.GetContainerClass(),
					Name:           "init",
					IsDynamic:      false,
					ReturnType:     nil,
					Signature:      "",
				},
				Args: nil,
				Base: v,
			}
			e1 := &stmt.InvokeStmt{
				Lvalue:     expr.Var{},
				Rvalue:     &invoke,
				Container:  visitor.Ir.GetContainerMethod(),
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e1)

			//if has __destruct, mock __destruct
			destruct := expr.InvokeInstanceExp{
				MethodRef: lang.MethodRef{
					DeclaringClass: visitor.Ir.GetContainerClass(),
					Name:           "__destruct",
					IsDynamic:      false,
					ReturnType:     nil,
					Signature:      "",
				},
				Args: nil,
				Base: v,
			}
			destructInvoke := &stmt.InvokeStmt{
				Lvalue:     expr.Var{},
				Rvalue:     &destruct,
				Container:  visitor.Ir.GetContainerMethod(),
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(destructInvoke)
		}
		if rvalue, ok := ret.(expr.FieldAccess); ok {
			e := &stmt.LoadFieldStmt{
				Lvalue:     v,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
		if rvalue, ok := visitor.Ctx.ret.(expr.ArrayAccess); ok {
			e := &stmt.LoadArrayStmt{
				Lvalue:     v,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
	} else if l, ok := visitor.Ctx.ret.(expr.FieldAccess); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.StoreFieldStmt{
				Lvalue:     l,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}

		if r, ok := visitor.Ctx.ret.(expr.FieldAccess); ok {
			temp := MakeTmpVar(n.Expr, nil)
			s2 := &stmt.LoadFieldStmt{
				Lvalue:     temp,
				Rvalue:     r,
				LineNumber: n.Expr.GetPosition().StartLine,
			}
			visitor.AddStmt(s2)

			s3 := &stmt.StoreFieldStmt{
				Lvalue:     l,
				Rvalue:     temp,
				LineNumber: n.Var.GetPosition().StartLine,
			}
			visitor.AddStmt(s3)
		}
	} else if l, ok := visitor.Ctx.ret.(expr.ArrayAccess); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.StoreArrayStmt{
				Lvalue:     l,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
	}
}

func (visitor *AstVisitor) AddStmt(stmt stmt.Stmt) {
	visitor.SeqIndex++
	stmt.SetStmtSeq(visitor.SeqIndex)
	visitor.Ir.AddStmt(stmt)
}

// this code block copied from ExprAssign
func (visitor *AstVisitor) ExprAssignReference(n *ast.ExprAssignReference) {
	visitor.drive1(n.Var, n, visitor.Ctx.parent)
	if v, ok := visitor.Ctx.ret.(expr.Var); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret

		if rvalue, ok := ret.(expr.BinaryExpr); ok {
			s := stmt.AssignBinaryStmt{
				DefinitionStmt: nil,
				Lvalue:         v,
				Rvalue:         rvalue,
				LineNumber:     n.Position.StartLine,
			}
			visitor.AddStmt(&s)
		}
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.AssignVarStmt{
				DefinitionStmt: nil,
				Lvalue:         v,
				Rvalue:         rvalue,
				LineNumber:     n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
		if rvalue, ok := ret.(expr.NewInstanceExpr); ok {
			e := &stmt.NewStmt{
				Lvalue:     v,
				Rvalue:     &rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
			//mock init
			invoke := expr.InvokeInstanceExp{
				MethodRef: lang.MethodRef{
					DeclaringClass: visitor.Ir.GetContainerClass(),
					Name:           "init",
					IsDynamic:      false,
					ReturnType:     nil,
					Signature:      "",
				},
				Args: nil,
				Base: v,
			}
			e1 := &stmt.InvokeStmt{
				Lvalue:     expr.Var{},
				Rvalue:     &invoke,
				Container:  visitor.Ir.GetContainerMethod(),
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e1)

			//if has __destruct, mock __destruct
			destruct := expr.InvokeInstanceExp{
				MethodRef: lang.MethodRef{
					DeclaringClass: visitor.Ir.GetContainerClass(),
					Name:           "__destruct",
					IsDynamic:      false,
					ReturnType:     nil,
					Signature:      "",
				},
				Args: nil,
				Base: v,
			}
			destructInvoke := &stmt.InvokeStmt{
				Lvalue:     expr.Var{},
				Rvalue:     &destruct,
				Container:  visitor.Ir.GetContainerMethod(),
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(destructInvoke)
		}
		if rvalue, ok := ret.(expr.FieldAccess); ok {
			e := &stmt.LoadFieldStmt{
				Lvalue:     v,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
		if rvalue, ok := visitor.Ctx.ret.(expr.ArrayAccess); ok {
			e := &stmt.LoadArrayStmt{
				Lvalue:     v,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
	} else if l, ok := visitor.Ctx.ret.(expr.FieldAccess); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.StoreFieldStmt{
				Lvalue:     l,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}

		if r, ok := visitor.Ctx.ret.(expr.FieldAccess); ok {
			temp := MakeTmpVar(n.Expr, nil)
			s2 := &stmt.LoadFieldStmt{
				Lvalue:     temp,
				Rvalue:     r,
				LineNumber: n.Expr.GetPosition().StartLine,
			}
			visitor.AddStmt(s2)

			s3 := &stmt.StoreFieldStmt{
				Lvalue:     l,
				Rvalue:     temp,
				LineNumber: n.Var.GetPosition().StartLine,
			}
			visitor.AddStmt(s3)
		}
	} else if l, ok := visitor.Ctx.ret.(expr.ArrayAccess); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.StoreArrayStmt{
				Lvalue:     l,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
	}
}

func (visitor *AstVisitor) ExprAssignBitwiseAnd(n *ast.ExprAssignBitwiseAnd) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignBitwiseOr(n *ast.ExprAssignBitwiseOr) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignBitwiseXor(n *ast.ExprAssignBitwiseXor) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignCoalesce(n *ast.ExprAssignCoalesce) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

// copy from ExprAssign
func (visitor *AstVisitor) ExprAssignConcat(n *ast.ExprAssignConcat) {
	visitor.drive1(n.Var, n, visitor.Ctx.parent)
	if v, ok := visitor.Ctx.ret.(expr.Var); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret

		if rvalue, ok := ret.(expr.BinaryExpr); ok {
			s := stmt.AssignBinaryStmt{
				DefinitionStmt: nil,
				Lvalue:         v,
				Rvalue:         rvalue,
				LineNumber:     n.Position.StartLine,
			}
			visitor.AddStmt(&s)
		}
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.AssignVarStmt{
				DefinitionStmt: nil,
				Lvalue:         v,
				Rvalue:         rvalue,
				LineNumber:     n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
		if rvalue, ok := ret.(expr.NewInstanceExpr); ok {
			e := &stmt.NewStmt{
				Lvalue:     v,
				Rvalue:     &rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
			//mock init
			invoke := expr.InvokeInstanceExp{
				MethodRef: lang.MethodRef{
					DeclaringClass: visitor.Ir.GetContainerClass(),
					Name:           "init",
					IsDynamic:      false,
					ReturnType:     nil,
					Signature:      "",
				},
				Args: nil,
				Base: v,
			}
			e1 := &stmt.InvokeStmt{
				Lvalue:     expr.Var{},
				Rvalue:     &invoke,
				Container:  visitor.Ir.GetContainerMethod(),
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e1)

			//if has __destruct, mock __destruct
			destruct := expr.InvokeInstanceExp{
				MethodRef: lang.MethodRef{
					DeclaringClass: visitor.Ir.GetContainerClass(),
					Name:           "__destruct",
					IsDynamic:      false,
					ReturnType:     nil,
					Signature:      "",
				},
				Args: nil,
				Base: v,
			}
			destructInvoke := &stmt.InvokeStmt{
				Lvalue:     expr.Var{},
				Rvalue:     &destruct,
				Container:  visitor.Ir.GetContainerMethod(),
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(destructInvoke)
		}
		if rvalue, ok := ret.(expr.FieldAccess); ok {
			e := &stmt.LoadFieldStmt{
				Lvalue:     v,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
		if rvalue, ok := visitor.Ctx.ret.(expr.ArrayAccess); ok {
			e := &stmt.LoadArrayStmt{
				Lvalue:     v,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
	} else if l, ok := visitor.Ctx.ret.(expr.FieldAccess); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.StoreFieldStmt{
				Lvalue:     l,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}

		if r, ok := visitor.Ctx.ret.(expr.FieldAccess); ok {
			temp := MakeTmpVar(n.Expr, nil)
			s2 := &stmt.LoadFieldStmt{
				Lvalue:     temp,
				Rvalue:     r,
				LineNumber: n.Expr.GetPosition().StartLine,
			}
			visitor.AddStmt(s2)

			s3 := &stmt.StoreFieldStmt{
				Lvalue:     l,
				Rvalue:     temp,
				LineNumber: n.Var.GetPosition().StartLine,
			}
			visitor.AddStmt(s3)
		}
	} else if l, ok := visitor.Ctx.ret.(expr.ArrayAccess); ok {
		visitor.drive1(n.Expr, n, visitor.Ctx.parent)
		ret := visitor.Ctx.ret
		if rvalue, ok := ret.(expr.Var); ok {
			e := &stmt.StoreArrayStmt{
				Lvalue:     l,
				Rvalue:     rvalue,
				LineNumber: n.Position.StartLine,
			}
			visitor.AddStmt(e)
		}
	}
}

func (visitor *AstVisitor) ExprAssignDiv(n *ast.ExprAssignDiv) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignMinus(n *ast.ExprAssignMinus) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignMod(n *ast.ExprAssignMod) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignMul(n *ast.ExprAssignMul) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignPlus(n *ast.ExprAssignPlus) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignPow(n *ast.ExprAssignPow) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignShiftLeft(n *ast.ExprAssignShiftLeft) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprAssignShiftRight(n *ast.ExprAssignShiftRight) {
	visitor.drive(n.Var)
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprBinaryBitwiseAnd(n *ast.ExprBinaryBitwiseAnd) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryBitwiseOr(n *ast.ExprBinaryBitwiseOr) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryBitwiseXor(n *ast.ExprBinaryBitwiseXor) {
	visitor.drive1(n.Left, n, visitor.Ctx.parent)
	l := visitor.Ctx.ret
	visitor.drive1(n.Right, n, visitor.Ctx.parent)
	r := visitor.Ctx.ret

	var b expr.BinaryExpr
	if v, ok := l.(expr.Var); ok {
		b = &expr.BitwiseExp{
			BinaryExpr: nil,
			Operator:   expr.XOR,
			Lop:        v,
			Rop:        r.(expr.Var),
		}
	}
	var v expr.Var
	if _, ok := visitor.Ctx.parent.(*ast.ExprAssign); ok {
		visitor.Ctx.ret = b
	} else {
		v = expr.Var{
			Rvalue: nil,
			Lvalue: nil,
			Name:   "%temp" + strconv.Itoa(n.Position.StartPos) + "_" + strconv.Itoa(n.Position.EndPos),
			Type:   &lang.UnKnowType{},
		}
		s := stmt.AssignBinaryStmt{
			DefinitionStmt: nil,
			Lvalue:         v,
			Rvalue:         b,
			LineNumber:     n.Position.StartLine,
		}
		visitor.AddStmt(&s)
		visitor.Ctx.ret = v
	}
}

func (visitor *AstVisitor) ExprBinaryBooleanAnd(n *ast.ExprBinaryBooleanAnd) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryBooleanOr(n *ast.ExprBinaryBooleanOr) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryCoalesce(n *ast.ExprBinaryCoalesce) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryConcat(n *ast.ExprBinaryConcat) {
	visitor.drive1(n.Left, n, visitor.Ctx.parent)
	l := visitor.Ctx.ret
	visitor.drive1(n.Right, n, visitor.Ctx.parent)
	r := visitor.Ctx.ret
	var rv expr.Var
	var lv expr.Var
	switch v := r.(type) {
	case expr.ArrayAccess:
		tmp := expr.Var{
			Name: "%temp" + strconv.Itoa(n.Right.GetPosition().StartPos) + "_" + strconv.Itoa(n.Right.GetPosition().EndPos),
			Type: nil,
		}
		load := &stmt.LoadArrayStmt{
			Lvalue:     tmp,
			Rvalue:     v,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(load)
		rv = tmp
	case expr.Var:
		rv = v
	default:
		if visitor.Debug {
			common.Log().Error("No match in ExprBinaryConcat")
		}
	}

	switch v := l.(type) {
	case expr.ArrayAccess:
		tmp := MakeTmpVar(n.Right, nil)
		load := &stmt.LoadArrayStmt{
			Lvalue:     tmp,
			Rvalue:     v,
			LineNumber: n.Position.StartLine,
		}
		visitor.AddStmt(load)
		lv = tmp
	case expr.Var:
		lv = v
	default:
		if visitor.Debug {
			common.Log().Error("No match in ExprBinaryConcat")
		}
	}

	var b expr.BinaryExpr
	b = &expr.StringConcatExp{
		BinaryExpr: nil,
		Operator:   expr.CONCAT,
		Lop:        lv,
		Rop:        rv,
	}
	var v expr.Var
	if _, ok := visitor.Ctx.parent.(*ast.ExprAssign); ok {
		visitor.Ctx.ret = b
	} else {
		v = MakeTmpVar(n, nil)
		s := stmt.AssignBinaryStmt{
			DefinitionStmt: nil,
			Lvalue:         v,
			Rvalue:         b,
			LineNumber:     n.Position.StartLine,
		}
		visitor.AddStmt(&s)
		visitor.Ctx.ret = v
	}
}

func (visitor *AstVisitor) ExprBinaryDiv(n *ast.ExprBinaryDiv) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryEqual(n *ast.ExprBinaryEqual) {
	visitor.drive1(n.Left, n, visitor.Ctx.parent)
	visitor.drive1(n.Right, n, visitor.Ctx.parent)
}

func (visitor *AstVisitor) ExprBinaryGreater(n *ast.ExprBinaryGreater) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryGreaterOrEqual(n *ast.ExprBinaryGreaterOrEqual) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryIdentical(n *ast.ExprBinaryIdentical) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryLogicalAnd(n *ast.ExprBinaryLogicalAnd) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryLogicalOr(n *ast.ExprBinaryLogicalOr) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryLogicalXor(n *ast.ExprBinaryLogicalXor) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryMinus(n *ast.ExprBinaryMinus) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryMod(n *ast.ExprBinaryMod) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryMul(n *ast.ExprBinaryMul) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryNotEqual(n *ast.ExprBinaryNotEqual) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryNotIdentical(n *ast.ExprBinaryNotIdentical) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryPlus(n *ast.ExprBinaryPlus) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryPow(n *ast.ExprBinaryPow) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryShiftLeft(n *ast.ExprBinaryShiftLeft) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinaryShiftRight(n *ast.ExprBinaryShiftRight) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinarySmaller(n *ast.ExprBinarySmaller) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinarySmallerOrEqual(n *ast.ExprBinarySmallerOrEqual) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprBinarySpaceship(n *ast.ExprBinarySpaceship) {
	visitor.drive(n.Left)
	visitor.drive(n.Right)
}

func (visitor *AstVisitor) ExprCastArray(n *ast.ExprCastArray) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprCastBool(n *ast.ExprCastBool) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprCastDouble(n *ast.ExprCastDouble) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprCastInt(n *ast.ExprCastInt) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprCastObject(n *ast.ExprCastObject) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprCastString(n *ast.ExprCastString) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ExprCastUnset(n *ast.ExprCastUnset) {
	visitor.drive(n.Expr)
}

func (visitor *AstVisitor) ScalarDnumber(n *ast.ScalarDnumber) {
	//todo
}

func (visitor *AstVisitor) ScalarEncapsed(n *ast.ScalarEncapsed) {
	var op expr.Var
	var ret expr.Var

	for i, part := range n.Parts {
		visitor.drive1(part, n, visitor.Ctx.parent)
		if i == 0 {
			op = visitor.Ctx.ret.(expr.Var)
			ret = op
		} else if i >= 1 {
			op = visitor.Ctx.ret.(expr.Var)
			tmp := MakeTmpVar(part, nil)
			rv := expr.StringConcatExp{
				Operator: expr.CONCAT,
				Lop:      ret,
				Rop:      op,
			}
			stmt := stmt.AssignBinaryStmt{
				Lvalue:     tmp,
				Rvalue:     &rv,
				LineNumber: n.Position.StartLine,
			}
			ret = tmp
			visitor.AddStmt(&stmt)
		}
	}
	visitor.Ctx.ret = ret
	//visitor.driveList(n.Parts)
}

func (visitor *AstVisitor) ScalarEncapsedStringPart(n *ast.ScalarEncapsedStringPart) {
	stmt := NewAssignStmtScalarString(string(n.Value), n)
	visitor.AddStmt(&stmt)
	visitor.Ctx.ret = stmt.Lvalue

	//visitor.Ctx.ret = expr.ScalarString{
	//	Value: string(n.Value),
	//}
}

func (visitor *AstVisitor) ScalarEncapsedStringVar(n *ast.ScalarEncapsedStringVar) {
	visitor.drive(n.Name)
	visitor.drive(n.Dim)
}

func (visitor *AstVisitor) ScalarEncapsedStringBrackets(n *ast.ScalarEncapsedStringBrackets) {
	visitor.drive(n.Var)
}

func (visitor *AstVisitor) ScalarHeredoc(n *ast.ScalarHeredoc) {
	visitor.driveList(n.Parts)
}

func (visitor *AstVisitor) ScalarLnumber(n *ast.ScalarLnumber) {
	lv := MakeConstVar(n, &lang.ScalarNumberType{})
	value, _ := strconv.Atoi(string(n.Value))
	rv := expr.ScalarNumber{
		Value: value,
	}
	s := &stmt.AssignLiteralStmt{
		Lvalue:     lv,
		Rvalue:     &rv,
		LineNumber: n.Position.StartLine,
	}
	visitor.Ctx.ret = lv
	visitor.AddStmt(s)
	visitor.VM.AddConst(s)
}

func (visitor *AstVisitor) ScalarMagicConstant(n *ast.ScalarMagicConstant) {
	value := string(n.Value)
	rv := &expr.ScalarString{
		Value: value,
	}
	tmp := MakeTmpVar(n, nil)
	stmt := &stmt.AssignLiteralStmt{
		Lvalue:     tmp,
		Rvalue:     rv,
		LineNumber: n.Position.StartLine,
	}
	visitor.AddStmt(stmt)
	visitor.Ctx.ret = tmp
}

func (visitor *AstVisitor) ScalarString(n *ast.ScalarString) {
	//value := string(n.Value)
	//stmt := NewAssignStmtScalarString(value, n)
	//visitor.Ctx.ret = stmt.Lvalue
	//visitor.AddStmt(&stmt)
	//visitor.VM.AddConst(&stmt)

	value := string(n.Value)
	rv := expr.ScalarString{
		Literal: nil,
		Value:   common.TrimString(value),
	}

	lv := MakeConstVar(n, &lang.ScalarStringType{})
	stmt := &stmt.AssignLiteralStmt{
		Lvalue:     lv,
		Rvalue:     &rv,
		LineNumber: n.Position.StartLine,
	}
	visitor.Ctx.ret = lv
	visitor.AddStmt(stmt)
	visitor.VM.AddConst(stmt)
}

func NewAssignStmtScalarString(data string, n ast.Vertex) stmt.AssignLiteralStmt {
	rv := expr.ScalarString{
		Literal: nil,
		Value:   common.TrimString(data),
	}

	lv := MakeConstVar(n, &lang.ScalarStringType{})
	stmt := stmt.AssignLiteralStmt{
		Lvalue:     lv,
		Rvalue:     &rv,
		LineNumber: n.GetPosition().StartLine,
	}
	return stmt
}

func (visitor *AstVisitor) NameName(n *ast.Name) {
	part0 := n.Parts[0].(*ast.NamePart)
	visitor.Ctx.ret = string(part0.Value)
}

func (visitor *AstVisitor) NameFullyQualified(n *ast.NameFullyQualified) {
}

func (visitor *AstVisitor) NameRelative(n *ast.NameRelative) {
}

func (visitor *AstVisitor) NameNamePart(n *ast.NamePart) {
}
