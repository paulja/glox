package main

import (
	"fmt"
	"glox"
)

func main() {
	var (
		l = glox.NewUnary(glox.NewToken(glox.TokenMinus, "-", nil, 1), glox.NewLiteral(123))
		o = glox.NewToken(glox.TokenStar, "*", nil, 1)
		r = glox.NewGrouping(glox.NewLiteral(45.67))
	)

	ex := glox.NewBinary(l, o, r)
	ap := glox.ASTPrinter{}

	fmt.Println(ap.Print(ex))
}
