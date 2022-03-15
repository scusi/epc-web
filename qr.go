package epc

import (
	"github.com/skip2/go-qrcode"
	"fmt"
)

// (epc *EPC) MarshalQR() - returns EPC as QR-code in PNG format
func (epc *EPC) MarshalQR() (qr []byte, err error) {
	qr, err = qrcode.Encode(fmt.Sprintf("%s", epc), qrcode.Medium, 256)
	if err != nil {
		return
	}
	return
}
