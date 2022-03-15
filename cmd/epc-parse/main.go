package main

import (
	"gitlab.scusi.io/flow/epc"
	"flag"
	"log"
	"fmt"
)

var (
	filename string
)

func init() {
	flag.StringVar(&filename, "f", "test.epc", "EPC file to parse")
}

func main() {
	flag.Parse()
	e, err := epc.ParseEPCFile(filename)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	fmt.Printf("%s", e)
}
