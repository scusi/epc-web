# epc-web

A simple web-application, written in golang, to generate a QR-Code that encodes a bank transaction.
Those QR-Codes are also called _Girocode_. Such a _Girocode_ can be scanned with mobile banking apps to initiate a bank transfer.
This way the user does not need to type all the bank transaction details, but can scan the QR-Code and directly authorize the transaction.

## Docker

Docker is the prefered way to use the web-app. The following sections describe how to create a docker image and how to run it.

### build a docker image

```
docker build -t epc-web:mybuild ./
```

### run epc-web docker container

```
docker run --rm -p 8080:80 epc-web:mybuild
```

Afterwards you can point your webbrowser to: `http://127.0.0.1:8080` to use epc-web.

![Screenshot_epc-web.png](/images/Screenshot_epc-web.png)
![Screenshot_epc-web_02.png](/images/Screenshot_epc-web_02.png)

### run with docker-compose

There is a `docker-compose.yml` file provided that can be used to run the web-app.
The docker-compose file is prepared to be used with _traefic_.

## run epc-web locally

Instead of useing docker you can also run the web-app locally. This is shown in the following example.

```
cd epc-web/
go run ./ -l :9999
```

After executing the above command point your browser to: `https://127.0.0.1:9999/`.


## Links

- Heise Artikel zum Thema: [Online-Banking: Rechnungen schneller mit QR-Codes überweisen](https://heise.de/-6543687)
- Beschreibung des EPC069-12 Standard: [EPC069-12](https://www.europeanpaymentscouncil.eu/sites/default/files/kb/file/2018-05/EPC069-12%20v2.1%20Quick%20Response%20Code%20-%20Guidelines%20to%20Enable%20the%20Data%20Capture%20for%20the%20Initiation%20of%20a%20SCT.pdf)

### Example Application Online

You can find an example application online at: [https://epc.scusi.io/](https://epc.scusi.io/).
