package epc

import (
	"text/template"
)

var t template.Template

type EPC_VERSION int
const (
	Version EPC_VERSION = iota
	V1	// 001
	V2	// 002
)

type EPC_ENCODING int
const (
	Encoding EPC_ENCODING = iota
	UTF8		// 1
	ISO88591	// 2
)

type EPC struct {
	Version		EPC_VERSION
	Encoding	EPC_ENCODING
	BIC		string		// size=11
	Name		string		// size=70
	IBAN		string		// size=34
	Amount		float
	SEPA_PURPOSE	int		// size=4
	SCR		string		// size=35
	SUBJECT		string		// size=140
	NOTE		string		// size=70
}

func New() (e *EPC) {
	e.Version = V1
	e.Encoding = UTF8
}
