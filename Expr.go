// Generated code, do not edit.

package glox

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

type Grouping struct {
	Expr

	Expression Expr
}

func newGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}

type Literal struct {
	Expr

	Value interface{}
}

func newLiteral(value interface{}) *Literal {
	return &Literal{Value: value}
}

type Unary struct {
	Expr

	Operator Token
	Right    Expr
}

func newUnary(operator Token, right Expr) *Unary {
	return &Unary{Operator: operator, Right: right}
}
