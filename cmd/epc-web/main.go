// skeleton for a API exposed on http
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"html/template"
	"log"
	"net/http"
	//"net/url"
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
	e, err := pD2epc(pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qrs, err := epc2b64QR(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		//r.ParseForm()
		pageData = urlparam2pD(r)
		e, err := pD2epc(pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qrs, err := epc2b64QR(e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageData["qrs"] = qrs
		err = t.ExecuteTemplate(w, "epcform", pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	case "POST":
		r.ParseForm()
		pageData["epcname"] = r.Form["epcname"][0]
		pageData["epciban"] = r.Form["epciban"][0]
		pageData["epcsubject"] = r.Form["epcsubject"][0]
		pageData["epcamount"] = r.Form["epcamount"][0]
		//pageData = urlparam2pD(r)
		up := pD2URLparam(pageData)
		pageData["epcurl"] = up
		
		e, err := pD2epc(pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qrs, err := epc2b64QR(e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageData["qrs"] = qrs
		//log.Printf("qrs = %s", qrs)
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
<html>
  <head>
    <title>EPC QR-Code Generator</title>
    <style>
    /* Style inputs, select elements and textareas */
	input[type=text], select, textarea{
  		width: 100%;
  		padding: 12px;
  		border: 1px solid #ccc;
  		border-radius: 4px;
  		box-sizing: border-box;
  		resize: vertical;
	}

/* Style the label to display next to the inputs */
label {
  padding: 12px 12px 12px 0;
  display: inline-block;
}

legend {
	font-size: large;
	margin-top: 25px;
}

.help-text {
	font-size: x-small;
}

/* Style the submit button */
input[type=submit] {
  background-color: #04AA6D;
  color: white;
  padding: 12px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  float: left;
  margin-top: 20px;
}

/* Style the reset button */
input[type=reset] {
  background-color: #FF0000;
  color: white;
  padding: 12px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  float: right;
  margin-top: 20px;
}

/* Style the container */
.container {
  border-radius: 5px;
  background-color: #f2f2f2;
  padding: 20px;
}

/* Floating column for labels: 25% width */
.col-25 {
  float: left;
  width: 25%;
  margin-top: 6px;
}

/* Floating column for inputs: 75% width */
.col-75 {
  float: left;
  width: 75%;
  margin-top: 6px;
}

/* Clear floats after the columns */
.row:after {
  content: "";
  display: table;
  clear: both;
}

/* Responsive layout - when the screen is less than 600px wide, make the two columns stack on top of each other instead of next to each other */
@media screen and (max-width: 600px) {
  .col-25, .col-75, input[type=submit] {
    width: 100%;
    margin-top: 0;
  }
} 
    </style>
  </head>
  <body>
  <div class="container"> <!-- start container -->
  <h1>EPC-QR-Code Generator</h1>
  <p>Mit Hilfe des folgenden Formulars kann ein EPC-QR-Code erstellt werden. EPC ist ein Standard um Bank-Überweisungen als QR-Code darzustellen. Diese QR-Codes kann man mit den meisten Online-Banking Apps scannen und spart sich so das lästige Eingeben der Überweisungsdetails.</p>
<p>Der EPC069-12 Standard findet sich <a href="https://www.europeanpaymentscouncil.eu/sites/default/files/kb/file/2018-05/EPC069-12%20v2.1%20Quick%20Response%20Code%20-%20Guidelines%20to%20Enable%20the%20Data%20Capture%20for%20the%20Initiation%20of%20a%20SCT.pdf">europeanpaymentcouncil.eu</a></p>
<form action="/" method="GET">
<fieldset>
<legend>Überweisungsempfänger</legend>
<div class="row">
  <div class="col-25">
  <label for="epcname">Name Kontoinhaber</label>
  </div>
  <div class="col-75">
  <input name="epcname" type="text" placeholder="Vorname Nachname" aria-describedby="epcNameHelpText" value="{{.epcname}}" required autofocus>
  <p class="help-text" id="epcNameHelpText">
	Der Name des Empängers darf höchstens 70 Zeichen lang sein. Erlaubt sind die Zeichen Buchstaben, Zahlen, Leerzeichen sowie die Zeichen /-?:().,+'.
  </p>
  </div>
</div>

<div class="row">
  <div class="col-25">
  	<label for="epciban">IBAN</label>
  </div>
  <div class="col-75">
  <input name="epciban" type="text" placeholder="DE..." value="{{.epciban}}" required>
  <p class="help-text" id="epcNameHelpText">
	Eine gültige IBAN, 34 Stellen lang.
  </p>
  </div>
</div>
</fieldset>

<fieldset>
<legend>Überweisungsdetails</legend>
<div class="row">
  <div class="col-25">
  	<label for="epcamount">Betrag in Euro</label>
  </div>
  <div class="col-75">
  	<input name="epcamount" type="text" placeholder="Betrag in EURO" value="{{.epcamount}}" required>
  	<p class="help-text" id="epcNameHelpText">Betrag in Euro als Fließzahl, das heißt mit Punkt statt Komma als Trenner zwischen Euro und Cent.</p>
</div>
</div>
<div class="row">
  <div class="col-25">
  	<label for="epcsubject">Verwendungszweck</label>
  </div>
  <div class="col-75">
  <input name="epcsubject" type="text" placeholder="Verwendungszweck" value="{{.epcsubject}}" required>
  <p class="help-text" id="epcNameHelpText">Ein Verwendungszweck oder eine Buchungsreferenz, maximal 140 Zeichen lang, Erlaubt sind Buchstaben, Zahlen, Leerzeichen sowie die Zeichen /-?:().,+'.</p>
</div>
</div>

<input type="submit"> <input type="reset">
</fieldset>
</form>
&nbsp;

<div class="row">
	<div class="col-25">
		<img name="epcqrcode" src="data:image/png;base64,{{.qrs}}" alt="QR-CODE" />
	</div>
	<div class="col-75">
		<p>Scanne den nebenstehenden QR-Code mit deiner Online-Banking App um die Überweisungsdaten zu übernehmen.</p>
		<pre>
Empfänger:        {{.epcname}}
IBAN:	          {{.epciban}}
Betrag:           {{.epcamount}} Euro
Verwendungszweck: {{.epcsubject}}
		</pre>
	</div>
</div>

<div>
<!-- URL: {{.epcurl}} -->
</div>
</div> <!-- end container -->
</body></html>
{{end}}
`
