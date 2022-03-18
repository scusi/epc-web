package epc

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"text/template"
)

var epctmpl = `{{define "EPC_Message"}}BCD
{{printf "%03d" .Version}}
{{.Encoding}}
SCT
{{.BIC}}
{{.Name}}
{{.IBAN}}
EUR{{.Amount}}
{{.SEPA_PURPOSE}}
{{.SCR}}
{{.SUBJECT}}
{{.NOTE}}{{end}}`

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
	ISO88592	// 3
	ISO88594	// 4
	ISO88595	// 5
	ISO88597	// 6
	ISO885910	// 7
	ISO885015	// 8
)

type EPC struct {
	Version		EPC_VERSION
	Encoding	EPC_ENCODING
	BIC		string		// size=11
	Name		string		// size=70
	IBAN		string		// size=34
	Amount		float64
	SEPA_PURPOSE	string		// size=4
	SCR		string		// size=35
	SUBJECT		string		// size=140
	NOTE		string		// size=70
}

func (epc *EPC) String() (s string) {
	//t, err := template.ParseFiles("EPC.tmpl")
	t, err := template.New("epc").Parse(epctmpl)
	if err != nil {
		log.Printf("eEPC.String() error parsing template: %s", err.Error())
		return s
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err = t.ExecuteTemplate(w, "EPC_Message", epc)
	if err != nil {
		log.Printf("EPC.String() error exec template: %s", err.Error())
		return s
	}
	w.Flush()
	return fmt.Sprintf("%s", b.Bytes())
}

func New(name, iban, subject string, ammount float64) (e *EPC) {
	e = new(EPC)
	e.Version = V2
	e.Encoding = UTF8
	e.Name = name
	e.IBAN = iban
	e.Amount = ammount
	e.SUBJECT = subject
	return e
}

func NewWithBIC(bic, name, iban, subject string, ammount float64) (e *EPC) {
	e = new(EPC)
	e.Version = V2
	e.Encoding = UTF8
	e.BIC = bic
	e.Name = name
	e.IBAN = iban
	e.Amount = ammount
	e.SUBJECT = subject
	return e
}
