FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/version"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
COPY core/types .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/version"
COPY exporter/version/go.mod .
COPY exporter/version/go.sum .

RUN go mod download

COPY exporter/version "/go/src/github.com/MaibornWolff/maDocK8s/exporter/version"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /version-exporter

FROM alpine:latest
RUN apk add --update ca-certificates

RUN mkdir -p "/etc/version/static"
COPY --from=builder /version-exporter /bin/version-exporter
COPY  exporter/version/static /etc/version/static
ENTRYPOINT ["/bin/version-exporter"]
