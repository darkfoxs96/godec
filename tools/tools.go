package tools

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func IncludeStr(arr []string, search string) bool {
	for _, v := range arr {
		if v == search {
			return true
		}
	}

	return false
}

func CopyFile(fromF, toF string) (err error) {
	from, err := os.Open(fromF)
	if err != nil {
		return
	}
	defer from.Close()

	to, err := os.OpenFile(toF, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return
	}

	return
}

type CMDWriter struct {
}

func (c *CMDWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return len(p), nil
}

func Command(name string, arg ...string) (err error) {
	cmd := exec.Command(name, arg...)
	var out CMDWriter
	var outErr CMDWriter

	cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &out
	cmd.Stderr = &outErr

	err = cmd.Run()
	if err != nil {
		fmt.Println("CMD ERR", err)
	}

	return
}

func BoolToStr(b bool, bTrue, bFalse string) string {
	if b {
		return bTrue
	} else {
		return bFalse
	}
}
