package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/paulja/glox"
)

var hadError bool

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(f string) {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		panic(fmt.Sprintf("glox: runFile error: %v", err))
	}
	run(string(buf))

	if hadError {
		os.Exit(2)
	}
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		s.Scan()
		run(s.Text())
		hadError = false
	}
}

func run(in string) {
	s := glox.NewScanner(in)
	go func() {
		for {
			select {
			case tkn := <-s.T:
				fmt.Println(tkn)
			case err := <-s.E:
				problem(err)
			case <-s.Done:
				break
			}
		}
	}()
	s.Scan()
}

func problem(err error) {
	switch v := err.(type) {
	case glox.LineError:
		fmt.Printf("line: %d, ", v.Line+1)
	}
	fmt.Printf("error: %s\n", err.Error())
	hadError = true
}
