package tmp

import (
	"fmt"
	"os"

	"github.com/darkfoxs96/godec/tools"
)

func GenTmp(ignore []string, vendor map[string]string, command string, isWin bool, isModify bool) (err error) {
	buildFile, err := os.Create("./godec_pack_b_b.go")
	if err != nil {
		return
	}
	defer buildFile.Close()

	imports := ""
	mapFunc := ""
	ignoreStr := ""
	for nameF, packF := range vendor {
		imports += `"` + packF + `"` + "\n"
		mapFunc += `"` + nameF + `": ` + nameF + ",\n"
	}

	for _, v := range ignore {
		ignoreStr += `"` + v + `",` + "\n"
	}

	fmt.Fprintln(buildFile, `package main

import(
	"github.com/darkfoxs96/godec/decbuilder"
	`+imports+`
)

var vendor_s = decbuilder.Vendors{
	`+mapFunc+`
}
var ignore_s = []string{
	`+ignoreStr+`
}

func main() {
	decbuilder.Builder(vendor_s, ignore_s, "`+command+`", `+tools.BoolToStr(isWin, "true", "false")+`, `+tools.BoolToStr(isModify, "true", "false")+`)
}
`)

	return
}
