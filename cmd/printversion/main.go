package main

import (
	"fmt"
	"os"
	"runtime/debug"

	_ "github.com/emirpasic/gods/lists/arraylist"
	_ "github.com/liyue201/gostl/ds/slice"
	_ "golang.org/x/exp/slices"
	_ "gonum.org/v1/gonum"
)

func main() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("failed")
		os.Exit(1)
	}

	for _, modDep := range buildInfo.Deps {
		fmt.Println("module name:", modDep.Path)
		fmt.Println("module version:", modDep.Version)
	}
}
