VERSION 0.7
FROM golang:1.21-alpine
WORKDIR /go-workdir

deps:
    # Download deps before copying code
    COPY go.mod go.sum .
    RUN go mod download
    # Output these back in case download changes hem
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum


build:
    FROM +deps
    COPY . .
    RUN go build -v -o output/server ./cmd/server/
    SAVE ARTIFACT output/server AS LOCAL local-output/server 

docker:
    ARG tag='latest'
    COPY .env .
    COPY +build/server .
    ENTRYPOINT ["/go-workdir/server"]
    SAVE IMAGE lenslocked-docker:$tag