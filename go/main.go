package main

import (
	"bufio"
	"fmt"
	"glox/interpreter"
	"glox/parser"
	"glox/resolver"
	"glox/scanner"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		log.Panic("cannot take more than one arguments")
		return
	} else if len(args) == 1 {
		runPrompt()
	} else {
		runFile(args[1])
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		if scanner.Text() == "exit" {
			os.Exit(64)
		}

		err := run(scanner.Text())
		if err != nil {
			fmt.Print(err)
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}
}

func runFile(file string) {
	prog, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}

	err = run(string(prog))
	if err != nil {
		fmt.Print(err)
		os.Exit(65)
	}
}

func run(source string) error {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		log.Panic(err)
	}

	parser := parser.NewParser(tokens)
	stmts, err := parser.Parse()
	if err != nil {
		return nil
	}

	interpreter := interpreter.NewInterpreter()
	resolver := resolver.NewResolver(interpreter)

	err = resolver.Resolve(stmts)
	if err != nil {
		return err
	}

	interpreter.Interpret(stmts)

	return nil
}
