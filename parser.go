package epc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	debug bool
)

func ParseEPCFile(filename string) (e *EPC, err error) {
	readFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	e = new(EPC)
	i := 1
	for fileScanner.Scan() {
		if i >= 12 {
			break
		}
		if debug {
			fmt.Printf("%02d\t'%s'\n", i, fileScanner.Text())
		}
		switch i {
		case 1:
			if fileScanner.Text() != "BCD" {
				err = fmt.Errorf("not a EPC file, BCD missing")
				return e, err
			}
		case 2:
			version, err := strconv.Atoi(fileScanner.Text())
			if err != nil {
				return e, err
			}
			e.Version = EPC_VERSION(version)
		case 3:
			encoding, err := strconv.Atoi(fileScanner.Text())
			if err != nil {
				return e, err
			}
			e.Encoding = EPC_ENCODING(encoding)
		case 4:
			if fileScanner.Text() != "SCT" {
				fmt.Errorf("not a EPC file, SCT missing")
				return e, err
			}
		case 5:
			e.BIC = fileScanner.Text()
		case 6:
			e.Name = fileScanner.Text()
		case 7:
			e.IBAN = fileScanner.Text()
		case 8:
			amountS := strings.TrimPrefix(fileScanner.Text(), "EUR")
			amount, err := strconv.ParseFloat(amountS, 64)
			if err != nil {
				return e, err
			}
			e.Amount = amount
		case 9:
			e.SEPA_PURPOSE = fileScanner.Text()
		case 10:
			e.SCR = fileScanner.Text()
		case 11:
			e.SUBJECT = fileScanner.Text()
		case 12:
			e.NOTE = fileScanner.Text()
		}
		i++
	}
	readFile.Close()
	return
}

/*
01	'BCD'
02	'002'
03	'1'
04	'SCT'
05	'COBADEFFXXX'
06	'Florian Walther'
07	'DE56120400000012262200'
08	'EUR123.42'
09	''
10	''
11	'Test Überweisung'

*/
