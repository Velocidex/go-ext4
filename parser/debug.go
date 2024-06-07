package parser

import "github.com/davecgh/go-spew/spew"

func Debug(arg interface{}) {
	spew.Dump(arg)
}
