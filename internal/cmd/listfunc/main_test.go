package main

import (
	"testing"

	. "github.com/spearson78/go-optic"
)

func TestFunctionsInFile(t *testing.T) {

	firstFunc, ok, err := GetFirst(functionsInFile(), "main.go")
	if !ok || err != nil || firstFunc.Name.Name != "main" {
		t.Fatal(firstFunc, ok, err)
	}

}
