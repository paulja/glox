package glox

import "fmt"

type ASTPrinter struct{}

func (a ASTPrinter) Print(ex Expr) string {
	return accept(a, ex)
}

func (a ASTPrinter) VisitBinaryExpr(ex *Binary) interface{} {
	return a.parenthesize(ex.Operator.Lexeme, ex.Left, ex.Right)
}

func (a ASTPrinter) VisitGroupingExpr(ex *Grouping) interface{} {
	return a.parenthesize("group", ex.Expression)
}

func (a ASTPrinter) VisitLiteralExpr(ex *Literal) interface{} {
	return fmt.Sprintf("%v", ex.Value)
}

func (a ASTPrinter) VisitUnaryExpr(ex *Unary) interface{} {
	return a.parenthesize(ex.Operator.Lexeme, ex.Right)
}

func (a ASTPrinter) parenthesize(name string, exprs ...Expr) string {
	var s string
	s += fmt.Sprintf("(%s", name)
	for _, ex := range exprs {
		s += fmt.Sprint(" ")
		s += accept(a, ex)
	}
	s += ")"
	return s
}

func accept(vtor Visitor, ex Expr) string {
	v := ex.Accept(vtor)
	switch v := v.(type) {
	case string:
		return v
	default:
		return ""
	}
}
