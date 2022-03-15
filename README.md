# go-epc

Sehr einfach implementierung des EPC069-12 Standard.

### Links

- Heise Artikel zum Thema: [Online-Banking: Rechnungen schneller mit QR-Codes überweisen](https://heise.de/-6543687)
- Beschreibung des EPC069-12 Standard: [EPC069-12](https://www.europeanpaymentscouncil.eu/sites/default/files/kb/file/2018-05/EPC069-12%20v2.1%20Quick%20Response%20Code%20-%20Guidelines%20to%20Enable%20the%20Data%20Capture%20for%20the%20Initiation%20of%20a%20SCT.pdf)

## download lib

```
go get gitlab.scusi.io/flow/epc
```

### build lib

```
cd $GOSRC/gitlab.scusi.io/flow/epc
go build ./
```

## build example program

```
cd $GOSRC/gitlab.scusi.io/flow/epc
go build ./cmd/epc-simple
```

## Usage

Mit dem folgenden beispiel Befehl unter linux kann man eine ECP Nachricht erzeugen.

```
$ ./epc-simple -i "DE53200400600200400600" -n "Bündnis Entwicklung Hilft" -a 5 -s "ARD/ Nothilfe Ukraine" -b "COBADEFFXXX" 
BCD
002
1
SCT
COBADEFFXXX
Bündnis Entwicklung Hilft
DE53200400600200400600
EUR5


ARD/ Nothilfe Ukraine
```

Wenn man den obigen Befehl an `qrencode -l H -t ANSIUTF8` piped kann man einen EPC QR-Code erzeugen.

```
$ ./epc-simple -i "DE53200400600200400600" -n "Bündnis Entwicklung Hilft" \
	-a 5 -s "ARD/ Nothilfe Ukraine" -b "COBADEFFXXX" \
	| qrencode -l H -t PNG -o images/test-qr.png 
```

Der obige befehl würde den folgenden QR-Code in die Datei `test-qr.png` schreiben.

![test-qr.png](/images/test-qr.png)


