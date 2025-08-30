package stmt

// StmtVisitor build boost components WorkList
type StmtVisitor interface {
	VisitAssignAssignBinaryStmt(stmt AssignBinaryStmt)
	VisitAssignLiteralStmt(stmt AssignLiteralStmt)
	VisitAssignVarStmt(stmt AssignVarStmt)
	VisitInvokeStaticStmt(stmt InvokeStmt)
	VisitNewStmt(stmt NewStmt)
	VisitLoadFieldStmt(stmt LoadFieldStmt)
	VisitStoreFieldStmt(stmt StoreFieldStmt)
	VisitStoreArrayStmt(stmt StoreArrayStmt)
	VisitLoadArrayStmt(stmt LoadArrayStmt)
}
