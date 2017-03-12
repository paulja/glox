// Generated code, do not edit.

package glox

type Visitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

type Expr interface {
	Accept(v Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewBinary(left Expr, operator *Token, right Expr) *Binary {
	return &Binary{Left: left, Operator: operator, Right: right}
}
func (n *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(n)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}
func (n *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(n)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{Value: value}
}
func (n *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(n)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func NewUnary(operator *Token, right Expr) *Unary {
	return &Unary{Operator: operator, Right: right}
}
func (n *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(n)
}
