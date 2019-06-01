package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/darkfoxs96/godec/tmp"
	"github.com/darkfoxs96/godec/tools"
)

type GlobalConfig struct {
	Command map[string]string
	Ignore  []string
	Vendor  map[string]string
}

func main() {
	testF := flag.Bool("test", false, "a dev mode. Don't use")
	modifyF := flag.Bool("modify", false, "not remove genareted files")
	runF := flag.String("run", "", "run your command script after building")
	comF := flag.String("com", "", "your command for app")
	goGetF := flag.Bool("get", false, "use go get ./...")
	flag.Parse()

	absDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	isWin := false
	if err != nil {
		panic(err)
	}

	if strings.Contains(absDir, "\\") { // this is windows os
		isWin = true
	}

	if *testF {
		absDir = "."
	}

	data, err := ioutil.ReadFile(absDir + "/godec.config.yaml")
	if err != nil {
		panic(err)
	}

	globalConfig := &GlobalConfig{}
	yaml.Unmarshal(data, globalConfig)

	c := ""
	if isWin {
		c = "go build -o project.exe"
	} else {
		c = "go build -o project"
	}

	if *runF != "" {
		c = globalConfig.Command[*runF]
		if c == "" {
			fmt.Println("Error: not found script '" + *runF + "'")
			return
		}

		c = strings.Replace(c, "{{.absDir}}", absDir, -1)
	}
	if *comF != "" {
		c += " " + *comF
	}
	c = strings.Replace(c, `\`, `\\`, -1)

	if err = tmp.GenTmp(globalConfig.Ignore, globalConfig.Vendor, c, isWin, *modifyF); err != nil {
		panic(err)
	}

	if *goGetF {
		if isWin {
			tools.Command("cmd", "/C", "cd "+absDir+" && go get ./...")
		} else {
			tools.Command("sh", "-c", "cd "+absDir+" && go get ./...")
		}
	}

	if isWin {
		tools.Command("cmd", "/C", "cd "+absDir+" && go run godec_pack_b_b.go")
	} else {
		tools.Command("sh", "-c", "cd "+absDir+" && go run godec_pack_b_b.go")
	}

	// if isWin {
	// 	tools.Command("cmd", "/C", "cd "+absDir+" && rd /s /q godec_pack_b_b")
	// } else {
	// 	tools.Command("sh", "-c", "cd "+absDir+" && rmdir --ignore-fail-on-non-empty godec_pack_b_b/")
	// }

	if err = os.Remove(absDir + "/godec_pack_b_b.go"); err != nil {
		panic(err)
	}

}
