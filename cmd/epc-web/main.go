// skeleton for a API exposed on http
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"encoding/base64"
	"gitlab.scusi.io/flow/epc"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var listenAddr string
var pageData = make(map[string]string)

func init() {
	flag.StringVar(&listenAddr, "l", "127.0.0.1:9999", "address to listen on, default is: 127.0.0.1:9999")
}

func main() {
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", EpcForm)
	log.Fatal(http.ListenAndServe(listenAddr, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Index, URL is: '%q'\n", html.EscapeString(r.URL.Path))
}


func EpcForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("epcform").Parse(epcformtmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":
		r.ParseForm()
		for k, v := range r.Form {
			pageData[k] = r.FormValue(v[0])
		}
		err = t.ExecuteTemplate(w, "epcform", pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	case "POST":
		r.ParseForm()
		amount, err := strconv.ParseFloat(r.Form["epc-amount"][0], 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := epc.New(
			r.Form["epc-name"][0],
			r.Form["epc-iban"][0],
			r.Form["epc-subject"][0],
			amount,
		)
		qr, err := e.MarshalQR()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}	
		qrs := base64.StdEncoding.EncodeToString(qr)
		pageData["qrs"] = qrs
		log.Printf("qrs = %s", qrs)
		err = t.ExecuteTemplate(w, "epcform", pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
		
	}
}

var epcformtmpl = `
{{define "epcform"}}
<html><head><title>EPC QR-Code Generator</title></head><body>

<form action="/" method="POST">
<fieldset>
<legend>Überweisungsempfänger</legend>
<div>
  <label for="epc-name">Name Kontoinhaber</label>
  <input name="epc-name" type="text" placeholder="Vorname Nachname" required autofocus>
</div>

<div>
  <label for="epc-iban">IBAN</label>
  <input name="epc-iban" type="text" placeholder="DE..." required>
</div>
</fieldset>

<fieldset>
<legend>Überweisungsdetails</legend>
<div>
  <label for="epc-amount">Betrag</label>
  <input name="epc-amount" type="text" placeholder="Betrag in EURO" required>
</div>
<div>
  <label for="epc-subject">Verwendungszweck</label>
  <input name="epc-subject" type="text" placeholder="Verwendungszweck" required>
</div>

<input type="submit"> <input type="reset">
</form>

<img src="data:image/png;base64,{{.qrs}}" alt="QR-CODE" />

</body></html>
{{end}}
`
