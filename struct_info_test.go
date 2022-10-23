package main

import (
	"fmt"
	goparser "go/parser"
	"go/token"
	"testing"
)

var codeStr = `
package main

// hello
type A struct {
}
`

func TestName(t *testing.T) {
	fset := token.NewFileSet()
	f, err := goparser.ParseFile(fset, "", []byte(codeStr), goparser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, c := range f.Comments {
		fmt.Println(c.Text())
	}
	fmt.Println(f.Doc)
	//for _, c := range f.Decls {
	//fmt.Println(c.(*ast.TypeSpec).Doc.Text())
	//}

}
