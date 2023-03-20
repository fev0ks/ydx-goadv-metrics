package checkers

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var CheckMainExit = &analysis.Analyzer{
	Name: "mainExitChecker",
	Run:  checkMainExit,
	Doc:  "Анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main.",
}

func checkMainExit(pass *analysis.Pass) (interface{}, error) {

	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isOsExitFuncCall(call) {
				pass.Reportf(x.Pos(), "os.Exit call is restricted in main func")
			}
		}
	}

	for _, file := range pass.Files {

		filename := pass.Fset.Position(file.Pos()).Filename
		if strings.HasSuffix(filename, "_test.go") || !strings.HasSuffix(filename, ".go") {
			continue
		}
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.String() != "main" {
					//fmt.Printf("'%s' - is not 'main' package\n", x.Name.String())
					return false
				}
			case *ast.FuncDecl:
				if x.Name.String() != "main" {
					//fmt.Printf("'%s' - is not 'main' function\n", x.Name.String())
					return false
				}
			case *ast.ExprStmt: // выражение
				expr(x)
			}
			return true
		})
	}
	return nil, nil
}

// *ast.ExprStmt {
//  X: *ast.CallExpr {
//    Fun: *ast.SelectorExpr {
//    .  X: *ast.Ident {
//    .  .  NamePos: 7:2
//    .  .  Name: "os"
//    .  }
//    .  Sel: *ast.Ident {
//    .  .  NamePos: 7:5
//    .  .  Name: "Exit"
//    .  }
//    }
//    Lparen: 7:9
//    Args: []ast.Expr (len = 1) {
//    .  0: *ast.BasicLit {
//    .  .  ValuePos: 7:10
//    .  .  Kind: INT
//    .  .  Value: "1"
//    .  }
//    }
//    Ellipsis: -
//    Rparen: 7:11

func isOsExitFuncCall(call *ast.CallExpr) bool {
	if callExpr, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := callExpr.X.(*ast.Ident); ok {
			if ident.Name == "os" && callExpr.Sel.Name == "Exit" {
				return true
			}
		}
	}
	return false
}
