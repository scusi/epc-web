// a simple webapp to generate EPC069-12 QR-Codes for bank transactions
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	branch    = "dev"
	version   = "v0.0.0"
	commit    = "000000"
	buildtime string
)

var debug bool
var listenAddr string
var pageData = make(map[string]string)
var t *template.Template
var err error

func init() {
	flag.StringVar(&listenAddr, "l", ":80", "address to listen on, default is: :80")
	flag.BoolVar(&debug, "debug", false, "logs debug info when set true")
}

func initTemplate() (t *template.Template, err error) {
	t, err = template.New("css").Parse(stylesheet)
	if err != nil {
		log.Printf("Error parsing stylesheet template: %s", err.Error())
		return
	}
	t, err = t.New("showQR").Parse(showQR)
	if err != nil {
		log.Printf("Error parsing showQR template: %s", err.Error())
		return
	}
	t, err = t.New("epcform").Parse(epcformtmpl)
	if err != nil {
		log.Printf("Error parsing EPCForm template: %s", err.Error())
		return
	}
	return
}

func main() {
	flag.Parse()
	if debug {
		log.Printf("Version: %s, Commit: %s, Branch: %s, Buildtime: %s", version, commit, branch, buildtime)
		log.Printf("listening on: %s", listenAddr)
	}
	mux := http.NewServeMux()
	mux.Handle("/", LogRequest(http.HandlerFunc(EpcForm)))
	mux.Handle("/version", LogRequest(http.HandlerFunc(Version)))
	mux.Handle("/qr", LogRequest(http.HandlerFunc(GetQR)))
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

func Version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Version: %s, Commit: %s, Branch: %s, Buildtime: %s\n", version, commit, branch, buildtime)
}

func GetQR(w http.ResponseWriter, r *http.Request) {
	t, err = initTemplate()
	if err != nil {
		log.Printf("Error parsing template: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	values := r.URL.Query()
	for k, v := range values {
		pageData[k] = v[0]
		if debug == true {
			log.Printf("read parameter %s : %s", k, v)
		}
	}
	e, err := pD2epc(pageData)
	if err != nil {
		log.Printf("Error creating EPC from pageData: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qrs, err := epc2b64QR(e)
	if err != nil {
		log.Printf("Error creating QR-Code from EPC data: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageData["qrs"] = qrs
	err = t.ExecuteTemplate(w, "showQR", pageData)
	if err != nil {
		log.Printf("Error executing template: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func EpcForm(w http.ResponseWriter, r *http.Request) {
	t, err = initTemplate()
	if err != nil {
		log.Printf("Error parsing template: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":
		//r.ParseForm()
		if len(r.URL.Query()) > 0 {
			pageData = urlparam2pD(r)
			e, err := pD2epc(pageData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			qrs, err := epc2b64QR(e)
			if err != nil {
				log.Printf("Error creating EPC from pageData: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pageData["qrs"] = qrs
		}
		err = t.ExecuteTemplate(w, "epcform", pageData)
		if err != nil {
			log.Printf("Error executing template: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	case "POST":
		r.ParseForm()
		pageData["epcname"] = r.Form["epcname"][0]
		pageData["epciban"] = r.Form["epciban"][0]
		pageData["epcbic"] = r.Form["epcbic"][0]
		pageData["epcsubject"] = r.Form["epcsubject"][0]
		pageData["epcamount"] = r.Form["epcamount"][0]
		//pageData = urlparam2pD(r)
		up := pD2URLparam(pageData)
		pageData["epcurl"] = up
		e, err := pD2epc(pageData)
		if err != nil {
			log.Printf("Error creating EPC from pageData: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qrs, err := epc2b64QR(e)
		if err != nil {
			log.Printf("Error creating QR-Code from EPC data: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageData["qrs"] = qrs
		if debug {
			log.Printf("qrs = %s", qrs)
		}
		err = t.ExecuteTemplate(w, "epcform", pageData)
		if err != nil {
			log.Printf("Error executing template: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

var stylesheet = `{{define "css"}}
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

pre.details {
	font-size: x-large;
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
{{end}}`

var epcformtmpl = `
{{define "epcform"}}
<html>
  <head>
    <title>EPC QR-Code Generator</title>
    <style>
    {{template "css"}}
    </style>
  </head>
  <body>
  <div class="container"> <!-- start container -->
  <h1>EPC-QR-Code Generator</h1>
  <p>Mit Hilfe des folgenden Formulars kann ein EPC-QR-Code erstellt werden. EPC ist ein Standard um Bank-Überweisungen als QR-Code darzustellen. Diese QR-Codes kann man mit den meisten Online-Banking Apps scannen und spart sich so das lästige Eingeben der Überweisungsdetails.</p>
<p>Der EPC069-12 Standard findet sich beim <a href="https://www.europeanpaymentscouncil.eu/sites/default/files/kb/file/2018-05/EPC069-12%20v2.1%20Quick%20Response%20Code%20-%20Guidelines%20to%20Enable%20the%20Data%20Capture%20for%20the%20Initiation%20of%20a%20SCT.pdf">europeanpaymentcouncil.eu</a></p>
<form action="/qr" method="GET">
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
  <p class="help-text" id="epcIBANHelpText">
	Eine gültige IBAN, 34 Stellen lang.
  </p>
  </div>
</div>
<div class="row">
  <div class="col-25">
  	<label for="epcbic">BIC</label>
  </div>
  <div class="col-75">
  <input name="epcbic" type="text" placeholder="BIC" value="{{.epcbic}}">
  <p class="help-text" id="epcBICHelpText">
	Eine gültige BIC. Optional, kann bei Inlandüberweisung weggelassen werden.
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


<div>
<!-- URL: {{.epcurl}} -->
</div>
</div> <!-- end container -->
</body></html>
{{end}}
`

var showQR = `{{define "showQR"}}
<html>
  <head>
    <title>EPC QR-Code Generator</title>
    <style>
    {{template "css"}}
    </style>
  </head>
  <body>
  <div class="container"> <!-- start container -->
  <h1>EPC-QR-Code Generator</h1>
  <ul>
    <li>
      <a href="/">Neuen QR-Code erstellen</a>
    </li>
    <li>
      <a href="/?epcname={{.epcname}}&epcbic={{.epcbic}}&epciban={{.epciban}}&epcamount={{.epcamount}}&epcsubject={{.epcsubject}}">Überweisungsdaten ändern</a>
    </li>
  </ul>
<div class="row">
	<div class="col-25">
		<img name="epcqrcode" src="data:image/png;base64,{{.qrs}}" alt="QR-CODE" />
	</div>
	<div class="col-75">
		<p>Scanne den nebenstehenden QR-Code mit deiner Online-Banking App um die Überweisungsdaten zu übernehmen.</p>
		<pre class="details">
Empfänger:        {{.epcname}}
IBAN:	          {{.epciban}}
{{if .epcbic -}}
BIC:		  {{.epcbic}}
{{end}}
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
