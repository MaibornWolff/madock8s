FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/gitlab"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
COPY core/types .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/utils"
COPY core/utils .

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/gitlab"
COPY exporter/gitlab/go.mod .
COPY exporter/gitlab/go.sum .

RUN go mod download

COPY exporter/gitlab "/go/src/github.com/MaibornWolff/maDocK8s/exporter/gitlab"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /gitlab

FROM alpine:latest
RUN apk add --update ca-certificates

RUN mkdir -p "/etc/gitlab/static"
COPY --from=builder /gitlab /bin/gitlab
COPY  exporter/gitlab/static /etc/gitlab/static
ENTRYPOINT ["/bin/gitlab"]
