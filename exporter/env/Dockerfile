FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/env"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
COPY core/types .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/utils"
COPY core/utils .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/env"
COPY exporter/env/go.mod .
COPY exporter/env/go.sum .

RUN go mod download

COPY exporter/env "/go/src/github.com/MaibornWolff/maDocK8s/exporter/env"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /env-exporter

FROM alpine:latest
RUN apk add --update ca-certificates

RUN mkdir -p "/etc/env-exporter/static"
COPY --from=builder /env-exporter /bin/env-exporter
COPY  exporter/env/static /etc/env-exporter/static
ENTRYPOINT ["/bin/env-exporter"]
