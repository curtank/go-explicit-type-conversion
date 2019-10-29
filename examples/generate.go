package main

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	// src is the input for which we want to print the AST.
	src := `
package main
func main() {
        println("Hello, World!")
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	fmt.Print(f)
	printer.Fprint(os.Stdout, fset, f)

}
