// Generated code, do not edit.

package glox

type visitor interface {
	visitBinaryExpr(expr *Binary) interface{}
	visitGroupingExpr(expr *Grouping) interface{}
	visitLiteralExpr(expr *Literal) interface{}
	visitUnaryExpr(expr *Unary) interface{}
}

type Expr struct{}

type Binary struct {
	Expr

	Left     Expr
	Operator Token
	Right    Expr
}

func newBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{Left: left, Operator: operator, Right: right}
}
func (n *Binary) accept(v visitor) interface{} {
	return v.visitBinaryExpr(n)
}

type Grouping struct {
	Expr

	Expression Expr
}

func newGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}
func (n *Grouping) accept(v visitor) interface{} {
	return v.visitGroupingExpr(n)
}

type Literal struct {
	Expr

	Value interface{}
}

func newLiteral(value interface{}) *Literal {
	return &Literal{Value: value}
}
func (n *Literal) accept(v visitor) interface{} {
	return v.visitLiteralExpr(n)
}

type Unary struct {
	Expr

	Operator Token
	Right    Expr
}

func newUnary(operator Token, right Expr) *Unary {
	return &Unary{Operator: operator, Right: right}
}
func (n *Unary) accept(v visitor) interface{} {
	return v.visitUnaryExpr(n)
}
