package decorators

import (
	"fmt"
	"go/ast"
	"os"
)

func GenMarshJSON(currType *ast.TypeSpec, g *ast.GenDecl, fName, pDir, packName string, params ...interface{}) (genFiles []string, err error) {
	genFiles = []string{pDir + "/" + currType.Name.Name + "_marsh_obj.go"}
	currentStrcut, ok := currType.Type.(*ast.StructType)
	if !ok {
		return
	}
	currentStrcut = currentStrcut

	fStr := `package ` + packName + `

func (a *` + currType.Name.Name + `) MarshalJSON() ([]byte, error) {
	return []byte(` + "`" + `{
		"text":"hihihhi"
	}` + "`" + `), nil
}
`

	buildFile, err := os.Create(pDir + "/" + currType.Name.Name + "_marsh_obj.go")
	if err != nil {
		return
	}
	defer buildFile.Close()
	fmt.Fprintln(buildFile, fStr)

	return
}
