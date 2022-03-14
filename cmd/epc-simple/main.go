package main

import (
	"gitlab.scusi.io/flow/epc"
	"flag"
	"fmt"
)

var bic		string
var iban	string
var name	string
var subject	string
var ammount	float64
var version     int
var encoding    int

func init() {
	flag.Float64Var(&ammount, "a", 123.42,                   "ammount to transfer")
	flag.StringVar(&bic,      "b", "COBADEFFXXX",            "BIC of the recipient")
	flag.StringVar(&iban,     "i", "DE56120400000012262200", "IBAN of the recipient")
	flag.StringVar(&name,     "n", "Florian Walther",        "Name of the recipient")
	flag.StringVar(&subject,  "s", "Test Überweisung",       "subject of the transfer")
	flag.IntVar(&version,     "v", 2,                        "version of EPC, can be 1 or 2")
	flag.IntVar(&encoding,    "e", 1,			 "encoding used, 1=UTF-8 2=ISO8859-1")
}

func main() {
	flag.Parse()
	e := epc.New(
		bic,
		name,
		iban,
		subject,
		ammount,
	)
	fmt.Printf("%s", e)
}
