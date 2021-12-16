
#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
ENV GO111MODULE=on

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/services/mdstorage"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/types"
RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/core/utils"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/"
COPY services/mdstorage/go.mod services/mdstorage
COPY services/mdstorage/go.sum services/mdstorage

COPY types/go.mod types/
COPY types/go.sum types/
COPY utils/go.mod utils/
COPY utils/go.sum utils/

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/core/services/mdstorage"
RUN go mod download

COPY . "/go/src/github.com/MaibornWolff/maDocK8s/core/"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /mdstorage

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /docs
COPY --from=builder /mdstorage /mdstorage
ENTRYPOINT ./mdstorage
LABEL Name=mdstorage Version=0.0.1
EXPOSE 3000
