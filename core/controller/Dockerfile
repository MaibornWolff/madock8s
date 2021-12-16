FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --update --no-cache git build-base

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/controller"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/"

COPY controller/go.mod controller
COPY controller/go.sum controller

COPY types/go.mod types
COPY types/go.sum types

WORKDIR /go/src/github.com/MaibornWolff/maDocK8s/core/controller
RUN go mod download

COPY . "/go/src/github.com/MaibornWolff/maDocK8s/core/"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /madock8s

FROM alpine:3.4
RUN apk add --update ca-certificates

COPY --from=builder /madock8s /bin/madock8s

ENTRYPOINT ["/bin/madock8s"]
