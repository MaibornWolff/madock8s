FROM golang:alpine AS generator
ENV GO111MODULE=on
RUN apk add --update --no-cache git

RUN mkdir -p "/go/src/github.com/MaibornWolff/maDocK8s/exporter/sample-metrics"

WORKDIR "/go/src/github.com/MaibornWolff/maDocK8s/exporter/sample-metrics"

RUN GO111MODULE=off go get -u "github.com/go-swagger/go-swagger/cmd/swagger"

COPY . .

RUN GO111MODULE=off swagger generate spec -o ./swagger.json --scan-models

FROM swaggerapi/swagger-ui:latest

COPY --from=generator "/go/src/github.com/MaibornWolff/maDocK8s/exporter/sample-metrics/swagger.json" /swagger.json

# ENV SWAGGER_JSON /swagger.json
# ENV API_KEY "**None**"
# ENV BASE_URL "/docs"
# ENV SWAGGER_JSON_URL ""

EXPOSE 8080

CMD ["sh", "/usr/share/nginx/run.sh"]
