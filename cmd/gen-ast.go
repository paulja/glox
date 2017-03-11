package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: gen-ast <output directory>")
		os.Exit(1)
	}

	out, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	defineAst(out, "Expr", []string{
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value interface{}",
		"Unary    : Operator Token, Right Expr",
	})
}

func defineAst(out, base string, types []string) {
	var src string

	src += fmt.Sprintln("// Generated code, do not edit.")
	src += fmt.Sprintln("")
	src += fmt.Sprintln("package glox")

	src += defineVisitor(base, types)

	// expr types
	src += fmt.Sprintln("")
	src += fmt.Sprintf("type %s struct {}", base)

	for _, t := range types {
		cls := strings.TrimRight(strings.Split(t, ":")[0], " ")
		fld := strings.TrimRight(strings.Split(t, ":")[1], " ")
		src += defineType(base, cls, fld)
	}

	path := fmt.Sprintf("%s/%s.go", out, base)
	if err := saveFile(path, src); err != nil {
		panic(err)
	}
}

func defineVisitor(base string, types []string) string {
	var src string

	// visitor interface
	src += fmt.Sprintln("")
	src += fmt.Sprintln("type visitor interface {")
	for _, t := range types {
		cls := strings.TrimRight(strings.Split(t, ":")[0], " ")
		src += fmt.Sprintf("visit%s%s(expr *%s) interface{}", cls, base, cls)
		src += fmt.Sprintln("")
	}
	src += fmt.Sprintln("}")

	return src
}

func defineType(base, cls, fld string) string {
	var src string

	src += fmt.Sprintln("")
	src += fmt.Sprintln("")
	src += fmt.Sprintf("type %s struct {", cls)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("%s", base)
	src += fmt.Sprintln("")
	src += fmt.Sprintln("")

	// fields
	fs := strings.Split(fld, ",")
	for _, f := range fs {
		src += fmt.Sprintln(f)
	}
	src += fmt.Sprintln("}")

	// new func
	src += fmt.Sprintf("func new%s(", cls)
	params := []string{}
	for _, f := range fs {
		t := strings.Split(f, " ")[2]
		n := strings.ToLower(strings.Split(f, " ")[1])
		params = append(params, fmt.Sprintf("%s %s", n, t))
	}
	src += fmt.Sprintf(strings.Join(params, ","))
	src += fmt.Sprintf(") *%s {", cls)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("return &%s{", cls)
	args := []string{}
	for _, f := range fs {
		t := strings.ToLower(strings.Split(f, " ")[1])
		n := strings.Split(f, " ")[1]
		args = append(args, fmt.Sprintf("%s: %s", n, t))
	}
	src += fmt.Sprintf(strings.Join(args, ","))
	src += fmt.Sprintln("}")
	src += fmt.Sprintln("}")

	// accept func
	src += fmt.Sprintf("func (n *%s) accept(v visitor) interface{} {", cls)
	src += fmt.Sprintln("")
	src += fmt.Sprintf("return v.visit%s%s(n)", cls, base)
	src += fmt.Sprintf("}")
	src += fmt.Sprintln("")

	return src
}

func saveFile(path, src string) error {
	// gofmt
	buf, err := format.Source([]byte(src))
	if err != nil {
		return err
	}
	// save
	ioutil.WriteFile(path, buf, 0644)

	return nil
}
