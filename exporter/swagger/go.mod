module github.com/MaibornWolff/maDocK8s/exporter/swagger

go 1.15

require (
	github.com/MaibornWolff/maDocK8s/core/types v0.0.0-00010101000000-000000000000
	github.com/MaibornWolff/maDocK8s/core/utils v0.0.0-00010101000000-000000000000
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/spec v0.19.9
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	google.golang.org/grpc v1.33.1
)

replace github.com/MaibornWolff/maDocK8s/core/types => ../../core/types

replace github.com/MaibornWolff/maDocK8s/core/utils => ../../core/utils
