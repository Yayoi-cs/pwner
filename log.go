package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"pwner/utils"
	"reflect"
	"runtime"
	"strings"
)

func Dump(vals ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println(vals...)
		return
	}

	varNames := extractVarNames(file, line)

	for i, v := range vals {
		name := "?"
		if i < len(varNames) {
			name = varNames[i]
		}

		fmt.Printf("[%s%s%s]: ", utils.Type2Color(name), name, utils.ColorReset)
		fmt.Printf("%s", utils.Type2Color(v))
		printHex(v)
		fmt.Print(utils.ColorReset)

		if i < len(vals)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
}

func printHex(v interface{}) {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num := val.Int()
		fmt.Printf("0x%x (%d)", num, num)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num := val.Uint()
		fmt.Printf("0x%x (%d)", num, num)

	case reflect.Slice, reflect.Array:
		fmt.Print("[")
		for i := 0; i < val.Len(); i++ {
			if i > 0 {
				fmt.Print(", ")
			}
			elem := val.Index(i)
			switch elem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fmt.Printf("0x%x", elem.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fmt.Printf("0x%02x", elem.Uint())
			default:
				fmt.Printf("%v", elem.Interface())
			}
		}
		fmt.Print("]")

	default:
		fmt.Printf("%+v", v)
	}
}

func extractVarNames(file string, line int) []string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	var varNames []string
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			pos := fset.Position(call.Pos())
			if pos.Line == line {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "Dump" {
					for _, arg := range call.Args {
						varNames = append(varNames, exprToString(arg))
					}
				}
			}
		}
		return true
	})

	return varNames
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.IndexExpr:
		return exprToString(e.X) + "[" + exprToString(e.Index) + "]"
	case *ast.CallExpr:
		return exprToString(e.Fun) + "()"
	case *ast.BasicLit:
		return e.Value
	default:
		return strings.TrimSpace(fmt.Sprintf("%T", expr))
	}
}
