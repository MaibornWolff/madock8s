FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/github"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
COPY core/types .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/utils"
COPY core/utils .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/github"
COPY exporter/github/go.mod .
COPY exporter/github/go.sum .

RUN go mod download

COPY exporter/github "/go/src/github.com/MaibornWolff/maDocK8s/exporter/github"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /github

FROM alpine:latest
RUN apk add --update ca-certificates

RUN mkdir -p "/etc/github/static"
COPY --from=builder /github /bin/github
COPY  exporter/github/static /etc/github/static
ENTRYPOINT ["/bin/github"]
