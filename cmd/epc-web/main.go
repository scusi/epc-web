// skeleton for a API exposed on http
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
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
	router.HandleFunc("/qr", GetQR)
	log.Fatal(http.ListenAndServe(listenAddr, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Index, URL is: '%q'\n", html.EscapeString(r.URL.Path))
}

func GetQR(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("epcform").Parse(epcformtmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	values := r.URL.Query()
	for k, v := range values {
		pageData[k] = v[0]
		log.Printf("read parameter %s : %s", k, v)
	}
	amount := 0.0
	if len(pageData["epcamount"]) > 0 {
		amount, err = strconv.ParseFloat(pageData["epcamount"], 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	e := epc.New(
		pageData["epcname"],
		pageData["epciban"],
		pageData["epcsubject"],
		amount,
	)
	qr, err := e.MarshalQR()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qrs := base64.StdEncoding.EncodeToString(qr)
	pageData["qrs"] = qrs
	err = t.ExecuteTemplate(w, "epcform", pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
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
		values := r.URL.Query()
		for k, v := range values {
			pageData[k] = v[0]
			log.Printf("read parameter %s : %s", k, v)
		}
		amount, err := strconv.ParseFloat(pageData["epcamount"], 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		e := epc.New(
			pageData["epcname"],
			pageData["epciban"],
			pageData["epcsubject"],
			amount,
		)
		qr, err := e.MarshalQR()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qrs := base64.StdEncoding.EncodeToString(qr)
		pageData["qrs"] = qrs
		err = t.ExecuteTemplate(w, "epcform", pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	case "POST":
		r.ParseForm()
		amount, err := strconv.ParseFloat(r.Form["epcamount"][0], 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := epc.New(
			r.Form["epcname"][0],
			r.Form["epciban"][0],
			r.Form["epcsubject"][0],
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
  <label for="epcname">Name Kontoinhaber</label>
  <input name="epcname" type="text" placeholder="Vorname Nachname" value="{{.epcname}}" required autofocus>
</div>

<div>
  <label for="epciban">IBAN</label>
  <input name="epciban" type="text" placeholder="DE..." value="{{.epciban}}" required>
</div>
</fieldset>

<fieldset>
<legend>Überweisungsdetails</legend>
<div>
  <label for="epcamount">Betrag</label>
  <input name="epcamount" type="text" placeholder="Betrag in EURO" value="{{.epcamount}}" required>
</div>
<div>
  <label for="epcsubject">Verwendungszweck</label>
  <input name="epcsubject" type="text" placeholder="Verwendungszweck" value="{{.epcsubject}}" required>
</div>

<input type="submit"> <input type="reset">
</form>

<img src="data:image/png;base64,{{.qrs}}" alt="QR-CODE" />

</body></html>
{{end}}
`
