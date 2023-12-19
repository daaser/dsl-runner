package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/daaser/dsl-runner/internal/dsl"
)

var passed, failed = 0, 0

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var directory string
	flag.StringVar(&directory, "d", "", "directory path containing test files")
	flag.Parse()

	if directory != "" {
		_, err := os.Stat(directory)
		if err != nil {
			log.Fatal(err)
		}
		if err := walkDir(directory); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := testRunner(os.Stdin); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("\nTESTS PASSED: %d\nTESTS FAILED: %d\n", passed, failed)

}

func walkDir(directory string) error {
	return filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		// log.Print(path)
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !d.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			return testRunner(f)
		}
		return nil
	})
}

func testRunner(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	if r == os.Stdin {
		fmt.Println("Enter assertions (Ctrl+D to exit):")
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Test: " + line)
		f, err := parser.ParseExpr(line)
		if err != nil {
			return err
		}
		expr, ok := f.(*ast.BinaryExpr)
		if !ok || (expr.Op != token.EQL && expr.Op != token.GTR && expr.Op != token.LSS && expr.Op != token.NEQ) {
			fmt.Println("only binary expressions are supported")
			failed++
			continue
		}
		res, err := dsl.EvalOp(expr)
		if err != nil {
			fmt.Println(err)
			failed++
			continue
		}
		fmt.Println(res)
		passed++
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
