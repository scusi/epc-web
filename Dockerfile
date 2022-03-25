## Dockerfile to create a docker image of epc-web application
#
# written by Florian Walther epc{at}scu.si
#
FROM golang:1.18.0-stretch AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1

WORKDIR /build

# Let's cache modules retrieval - those don't change so often
COPY go.mod .
#COPY go.sum .
RUN go mod download

# Copy the code necessary to build the application
# You may want to change this to copy only what you actually need.
COPY . .

# Build the application
#RUN go build ./cmd/HelloDocker
RUN go build ./cmd/epc-web

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
 VERSION=$(git describe --tags) && \
 BUILDTIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ') && \
 BRANCH=$(git branch | grep \* | cut -d ' ' -f2) && \
 go build -ldflags "-s -w -X main.commit=$GIT_COMMIT -X main.version=$VERSION -X main.branch=$BRANCH -X main.buildtime=$BUILDTIME" ./cmd/epc-web




# Let's create a /dist folder containing just the files necessary for runtime.
# Later, it will be copied as the / (root) of the output image.
WORKDIR /dist
RUN cp /build/epc-web ./epc-web

# Optional: in case your application uses dynamic linking (often the case with CGO), 
# this will collect dependent libraries so they're later copied to the final image
# NOTE: make sure you honor the license terms of the libraries you copy and distribute
RUN ldd epc-web | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

# Copy or create other directories/files your app needs during runtime.
# E.g. this example uses /data as a working directory that would probably
#      be bound to a perstistent dir when running the container normally
RUN mkdir /data

# Create the minimal runtime image
FROM scratch

COPY --chown=0:0 --from=builder /dist /
COPY --chown=0:0 --from=builder /build/epc-web /
#COPY --chown=0:0 config.yml /data/

# Set up the app to run as a non-root user inside the /data folder
# User ID 65534 is usually user 'nobody'. 
# The executor of this image should still specify a user during setup.
COPY --chown=65534:0 --from=builder /data /data
USER 65534
WORKDIR /data

ENTRYPOINT ["/epc-web"]
