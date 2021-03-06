package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chuongtrh/godepviz/godep"
)

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Need package name to process")
	}
	pkgName := args[0]
	node := &godep.Node{
		PkgName: pkgName,
		IsRoot:  true,
		Parent:  nil,
	}
	err := node.FindImports()
	if err != nil {
		log.Fatal(err)
	}
	graph := node.BuildGraph()
	fmt.Println(graph)
}
