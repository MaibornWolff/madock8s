FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/prometheus"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
COPY core/types .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/utils"
COPY core/utils .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/prometheus"
COPY exporter/prometheus/go.mod .
COPY exporter/prometheus/go.sum .

RUN go mod download

COPY exporter/prometheus "/go/src/github.com/MaibornWolff/maDocK8s/exporter/prometheus"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /prometheus

FROM alpine:latest
RUN apk add --update ca-certificates

RUN mkdir -p "/etc/prometheus/static"
COPY --from=builder /prometheus /bin/prometheus
COPY  exporter/prometheus/static /etc/prometheus/static
ENTRYPOINT ["/bin/prometheus"]