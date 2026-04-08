# STAGE 1: Build
FROM golang:1.22-alpine AS builder

# Build-Argumente für die Versionierung
ARG VERSION
ARG GIT_COMMIT
ARG BUILDTIME
ARG BRANCH

# Git wird nur für die Go-Module benötigt, falls diese aus privaten Repos kommen
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Abhängigkeiten cachen
COPY go.mod ./
# COPY go.sum ./  # Falls vorhanden
RUN go mod download

# Quellcode kopieren
COPY . .

# Statischer Build: CGO_ENABLED=0 ist zwingend für scratch!
# Wir nutzen die übergebenen ARGs in den ldflags
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w \
    -X main.version=${VERSION} \
    -X main.commit=${GIT_COMMIT} \
    -X main.buildtime=${BUILDTIME} \
    -X main.branch=${BRANCH}" \
    -o epc-web .

# STAGE 2: Final (Minimalistisches Image)
FROM scratch

# Root-Zertifikate für ausgehende HTTPS-Verbindungen
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Das Binary aus der Builder-Stage kopieren
COPY --from=builder /app/epc-web /epc-web

# Standardmäßig auf Port 8080 (anpassen falls deine App einen anderen nutzt)
EXPOSE 8080

ENTRYPOINT ["/epc-web"]
