FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/swagger"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
COPY core/types .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/utils"
COPY core/utils .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/swagger"
COPY exporter/swagger/go.mod .
COPY exporter/swagger/go.sum .

RUN go mod download

COPY exporter/swagger "/go/src/github.com/MaibornWolff/maDocK8s/exporter/swagger"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /swagger

FROM alpine:latest
RUN apk add --update ca-certificates

RUN mkdir -p "/etc/swagger/static"
COPY --from=builder /swagger /bin/swagger
COPY  exporter/swagger/static /etc/swagger/static
ENTRYPOINT ["/bin/swagger"]
