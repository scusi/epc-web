Versuch EPC QR-Codes in golang zu implementieren

Heise Artikel zum Thema: [Online-Banking: Rechnungen schneller mit QR-Codes überweisen](https://heise.de/-6543687)



## Usage

Mit dem folgenden beispiel befehl unter linux kann man eine ECP Nachricht erzeugen.

```
$ ./epc-simple -v 2 -e 1 -n "Sylvester Stallone" -s "Ein toller Test, für Döner" -a 23.42 -i DE56120400000012262200 -b COBADEFFXXX
BCD
002
1
SCT
COBADEFFXXX
Sylvester Stallone
DE56120400000012262200
EUR23.42


Ein toller Test, für Döner

```

Wenn man den obigen Befehl an `qrencode -l H -t ANSIUTF8` piped kann man einen EPC QR-Code erzeugen.

```
$ ./epc-simple -v 2 -e 1 -n "Sylvester Stallone" \
	-s "Ein toller Test, für Döner" -a 23.42 \
	-i DE56120400000012262200 -b COBADEFFXXX \
	| qrencode -l H -t PNG > test-qr.png
```

Der obige befehl würde den folgenden QR-Code in die Datei `test-qr.png` schreiben.

![test-qr.png](/images/test-qr.png)


