package parser

import (
	"github.com/davecgh/go-spew/spew"
)

var (
	debug bool
)

func Debug(arg interface{}) {
	spew.Dump(arg)
}

func DebugPrint(fmt_str string, v ...interface{}) {
	// fmt.Printf(fmt_str, v...)
}
