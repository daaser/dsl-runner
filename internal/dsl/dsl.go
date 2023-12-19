package dsl

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func Eval(expr ast.Expr) (int, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		return EvalBinary(e)
	case *ast.BasicLit:
		return EvalLit(e)
	}
	return -1, fmt.Errorf("illegal expression \"%s\"", reflect.TypeOf(expr))
}

func EvalLit(lit *ast.BasicLit) (int, error) {
	if lit.Kind != token.INT {
		return -1, fmt.Errorf("bad AST literal \"%s\"; only integers are supported", lit.Kind)
	}
	return strconv.Atoi(lit.Value)
}

func EvalOp(b *ast.BinaryExpr) (bool, error) {
	x, err := Eval(b.X)
	if err != nil {
		return false, err
	}
	y, err := Eval(b.Y)
	if err != nil {
		return false, err
	}

	switch b.Op {
	case token.EQL:
		return x == y, nil
	case token.GTR:
		return x > y, nil
	case token.LSS:
		return x < y, nil
	case token.NEQ:
		return x != y, nil
	}
	return false, fmt.Errorf("\"%s\" is not a valid token; allowed is ==, >, <, !=", b.Op)
}

func EvalBinary(b *ast.BinaryExpr) (int, error) {
	x, err := Eval(b.X)
	if err != nil {
		return -1, err
	}
	y, err := Eval(b.Y)
	if err != nil {
		return -1, err
	}

	switch b.Op {
	case token.ADD:
		return x + y, nil
	case token.SUB:
		return x - y, nil
	case token.MUL:
		return x * y, nil
	case token.QUO:
		return x / y, nil
	}
	return -1, fmt.Errorf("\"%s\" is not a valid token; allowed is +, -, *, /", b.Op)
}
