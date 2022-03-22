package main

import (
	"encoding/base64"
	"gitlab.scusi.io/flow/epc"
	"strconv"
	"log"
	"net/http"
	"net/url"
)

func urlparam2pD(r *http.Request) (pageData map[string]string) {
	pageData = make(map[string]string)
	values := r.URL.Query()
	for k, v := range values {
		pageData[k] = v[0]
		if debug { log.Printf("read parameter %s : %s", k, v)}
	}
	up := pD2URLparam(pageData)
	pageData["epcurl"] = up
	return	
}

func pD2URLparam(pageData map[string]string) (up string) {
	// encode URL values
	uv := url.Values{}
	for k, v := range pageData {
		uv.Add(k, v)
	}
	up = uv.Encode()
	return
}

func pD2epc(pageData map[string]string) (e *epc.EPC, err error) {
	amount := 0.0
	if len(pageData["epcamount"]) > 0 {
		amount, err = strconv.ParseFloat(pageData["epcamount"], 64)
		if err != nil {
			return
		}
	}
	e = epc.New(
		pageData["epcname"],
		pageData["epciban"],
		pageData["epcsubject"],
		amount,
	)
	return
}

func epc2b64QR(e *epc.EPC) (b64QR string, err error) {
	qr, err := e.MarshalQR()
	if err != nil {
		return
	}
	b64QR = base64.StdEncoding.EncodeToString(qr)
	return
}
