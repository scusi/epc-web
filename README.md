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
go build ./cmd/epc-parse
```

## Usage

### Text Format
Mit dem folgenden beispiel Befehl unter linux kann man eine ECP Nachricht erzeugen.

```
$ ./epc-simple -i "DE53200400600200400600" -n "Bündnis Entwicklung Hilft" \
	-a 5 -s "ARD/ Nothilfe Ukraine" -b "COBADEFFXXX" -format text 
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
	-a 5 -s "ARD/ Nothilfe Ukraine" -b "COBADEFFXXX" -format text\
	| qrencode -l H -t PNG -o images/test-qr.png 
```

Der obige befehl würde den folgenden QR-Code in die Datei `images/test-qr.png` schreiben.

![test-qr.png](/images/test-qr.png)

### PNG Format

`epc` kann auch selber QR-codes erzeugen, dazu muss man das format auf `png` setzten, etwa so:
```
$ ./epc-simple -i "DE53200400600200400600" -n "Bündnis Entwicklung Hilft" \
	-a 5 -s "ARD/ Nothilfe Ukraine" -b "COBADEFFXXX" -format png > images/test-qr2.png 
```

Der erzeugte QR-Code findet sich dann in der Daei `images/test-qr2.png`.

![test-qr2.png](/images/test-qr2.png)

### Parsing EPC Messages

```
$ ./epc-simple -i "DE53200400600200400600" -n "Bündnis Entwicklung Hilft" \
        -a 5 -s "ARD/ Nothilfe Ukraine" -b "COBADEFFXXX" -format text > test.epc
$ ./epc-parse -f test.epc
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

The first command from the above example does create a new EPC Message, in Text format and pipes that to a file called `test.epc`
The second command reads `test.epc` and parses the content into a EPC datastruct (`epc.EPC`), before writeing it to STDOUT.


## epc-web example application

There is an epc-web example application.

### build epc-web docker container

```
docker build -t epc-web:latest ./
```

### run epc-web docker container

```
docker run --rm -p 8080:80 epc-web:latest
```

Afterwards you can point your webbrowser to: [http://127.0.0.1:8080](http://127.0.0.1:8080/) to use epc-web.

![Screenshot_epc-web.png](/images/Sreenshot_epc-web.png)
