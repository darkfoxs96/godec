package decbuilder

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/darkfoxs96/godec/tools"
)

const PREFIX_PATH = "./"

type Vendors map[string]func(currType *ast.TypeSpec, g *ast.GenDecl, fName, pDir, packName string, params ...interface{}) (genFiles []string, err error)

var (
	vendor      = Vendors{}
	ignore      = []string{}
	command     string
	isWin       bool
	removeFilse = []string{}
)

func Builder(ven Vendors, ign []string, comm string, isW, isModify bool) {
	vendor, ignore, command, isWin = ven, ign, comm, isW

	pDir := ""
	if err := parseDir(pDir); err != nil {
		panic(err)
	}

	if command != "" {
		if isWin {
			tools.Command("cmd", "/C", command)
		} else {
			tools.Command("sh", "-c", command)
		}
	}

	if !isModify {
		for _, vF := range removeFilse {
			if err := os.Remove(vF); err != nil {
				panic(err)
			}
		}
	}
}

func parseDir(pDir string) (err error) {
	if pDir == "" {
		pDir = "./"
	} else {
		pDir += "/"
	}

	err = os.MkdirAll(PREFIX_PATH+pDir[2:], os.ModePerm)
	if err != nil {
		return
	}

	dirs := []string{}
	files, err := ioutil.ReadDir(pDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "godec_pack_b_b") || tools.IncludeStr(ignore, pDir+file.Name()) {
			continue
		}

		if file.IsDir() {
			dirs = append(dirs, pDir+file.Name())
		} else {
			// pathOrg := pDir + file.Name()
			// pathTo := PREFIX_PATH + pathOrg[2:]
			// err = tools.CopyFile(pathOrg, pathTo)
			// if err != nil {
			// 	return
			// }
		}
	}

	fset := token.NewFileSet()
	pack, err := parser.ParseDir(fset, pDir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, val := range pack {
		for fName, fileP := range val.Files {
			//
			for _, val := range fileP.Decls {
				g, ok := val.(*ast.GenDecl)
				if !ok {
					continue
				}

				for _, spec := range g.Specs {
					lI := strings.LastIndex(fName, "/")
					if lI == -1 {
						lI = strings.LastIndex(fName, "\\")
					}

					parseGenDecl(spec, g, fName[lI+1:], PREFIX_PATH+pDir[2:], fileP.Name.Name)
				}
			}
			//
		}
	}

	for _, d := range dirs {
		err = parseDir(d)
		if err != nil {
			return
		}
	}

	return
}

func parseGenDecl(spec ast.Spec, g *ast.GenDecl, fName, pDir, packName string) (err error) {
	currType, ok := spec.(*ast.TypeSpec)
	currType = currType
	if !ok {
		return
	}

	if g.Doc == nil {
		return
	}

	for _, comment := range g.Doc.List {
		posStart := strings.LastIndex(comment.Text, "//@:")
		if posStart == -1 {
			posStart = strings.LastIndex(comment.Text, "// @:")
			if posStart != -1 {
				posStart++
			}
		}

		if posStart == -1 {
			continue
		}

		posStart += 4
		startFnText := comment.Text[posStart:]
		posEndFn := strings.LastIndex(startFnText, "(")
		fnName := startFnText[:posEndFn]
		ven := vendor[fnName]
		if ven == nil {
			fmt.Println("Warning: not found '" + fnName + "' vendor!")
			continue
		}

		newFilse, err := ven(currType, g, fName, pDir, packName)
		if err != nil {
			return err
		}

		removeFilse = append(removeFilse, newFilse...)
	}

	return
}
