package epc

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"text/template"
)

// epctmpl - variable that contains the template for an EPC message
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

// t - variable that holds the parsed template during runtime
var t template.Template

// EPC_VERSION - a custom type for the EPC version
type EPC_VERSION int

const (
	Version EPC_VERSION = iota
	V1                  // 001
	V2                  // 002
)

// EPC_ENCODING a data type for the used encoding
type EPC_ENCODING int

const (
	Encoding  EPC_ENCODING = iota
	UTF8                   // 1
	ISO88591               // 2
	ISO88592               // 3
	ISO88594               // 4
	ISO88595               // 5
	ISO88597               // 6
	ISO885910              // 7
	ISO885015              // 8
)

// EPC - datastruct of an EPC message
type EPC struct {
	Version      EPC_VERSION
	Encoding     EPC_ENCODING
	BIC          string // size=11
	Name         string // size=70
	IBAN         string // size=34
	Amount       float64
	SEPA_PURPOSE string // size=4
	SCR          string // size=35
	SUBJECT      string // size=140
	NOTE         string // size=70
}

// (epc *EPC) String() - Stringer function for EPC datastructure
func (epc *EPC) String() (s string) {
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

// New(name, IBAN, subject string, ammount float64)
// returns a pointer to a new EPC datastructure and error, if any.
func New(name, iban, subject string, ammount float64) (e *EPC, err error) {
	if nameIsValid(name) != true {
		err = fmt.Errorf("supplied name is not valid")
		return
	}
	if ibanIsValid(iban) != true {
		err = fmt.Errorf("supplied IBAN is not valid")
		return
	}
	if subjectIsValid(subject) != true {
		err = fmt.Errorf("supplied subject is not valid")
		return
	}
	e = new(EPC)
	e.Version = V2
	e.Encoding = UTF8
	e.Name = name
	e.IBAN = iban
	e.Amount = ammount
	e.SUBJECT = subject
	return
}

// New(BIC, name, IBAN, subject string, ammount float64)
// returns a pointer to a new EPC datastructure and error, if any.
func NewWithBIC(bic, name, iban, subject string, ammount float64) (e *EPC, err error) {
	if BICIsValid(bic) != true {
		err = fmt.Errorf("supplied BIC is not valid")
		return
	}
	if nameIsValid(name) != true {
		err = fmt.Errorf("supplied name is not valid")
		return
	}
	if ibanIsValid(iban) != true {
		err = fmt.Errorf("supplied IBAN is not valid")
		return
	}
	if subjectIsValid(subject) != true {
		err = fmt.Errorf("supplied subject is not valid")
		return
	}
	e = new(EPC)
	e.Version = V2
	e.Encoding = UTF8
	e.BIC = bic
	e.Name = name
	e.IBAN = iban
	e.Amount = ammount
	e.SUBJECT = subject
	return
}

// ibanIsValid(IBAN) - helper function to validate a IBAN
// returns `true` if given string is a syntactical valid IBAN
func ibanIsValid(iban string) bool {
	ibanRE := regexp.MustCompile(`^[A-Z]{2}[0-9A-Z]{2,32}$`)
	return ibanRE.MatchString(iban)
}

// nameIsValid(Name) - helper function to validate a given Name
// returns `true` if given string is a syntactical valid Name
func nameIsValid(name string) bool {
	nameRE := regexp.MustCompile(`^[0-9A-Za-z\s\?\-\:\(\)\.\,\+\'üäöÜÄÖß]{2,70}$`)
	return nameRE.MatchString(name)
}

// subjectIsValid(subject) - helper function to validate a given Subject
// returns `true` if given string is a syntactical valid subject
func subjectIsValid(subject string) bool {
	subjectRE := regexp.MustCompile(`^[0-9A-Za-z\s\?\-\:\(\)\.\,\+\'üäöÜÄÖß]{1,140}$`)
	return subjectRE.MatchString(subject)
}

// BICIsValid(BIC) - helper function to validate a BIC (Bank Identification Code)
// returns `true` if given string is a syntactical valid BIC
func BICIsValid(bic string) bool {
	bicRE := regexp.MustCompile(`(?m)([a-zA-Z]{4})([a-zA-Z]{2})(([2-9a-zA-Z]{1})([0-9a-np-zA-NP-Z]{1}))((([0-9a-wy-zA-WY-Z]{1})([0-9a-zA-Z]{2}))|([xX]{3})|)`)
	return bicRE.MatchString(bic)
}
