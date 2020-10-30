FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/prometheus/sample-metrics"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/prometheus/sample-metrics"

COPY go.mod .
COPY go.mod .
RUN go mod download

COPY . "/go/src/github.com/MaibornWolff/maDocK8s/exporter/prometheus/sample-metrics"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /gometrics

FROM alpine:latest
RUN apk add --update ca-certificates

ENV CUSTOM_ENVIRONMENT_VARIABLE_FROM_DOCKERFILE "Hello, World!"

COPY --from=builder /gometrics /bin/gometrics

ENTRYPOINT ["/bin/gometrics"]