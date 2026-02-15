package main

import (
	"flag"
)

type multiFlag []string

func (i *multiFlag) String() string {
	return "multiFlag"
}

func (i *multiFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var imports multiFlag
var primitives multiFlag
var packagePrefix = flag.String("prefix", "", "package prefix")
var sourceConstraint = flag.String("constraint", "", "constraint foe the source type")

func parseArgs() {
	flag.Var(&imports, "import", "additional imports")
	flag.Var(&primitives, "primitive", "additional primitives")
	flag.Parse()
}

var rootObjName = flag.String("root", "DB", "root object name")
