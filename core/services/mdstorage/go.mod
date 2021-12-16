module github.com/MaibornWolff/maDocK8s/core/services/mdstorage

go 1.15

require (
	github.com/MaibornWolff/maDocK8s/core/types v0.0.0
	github.com/MaibornWolff/maDocK8s/core/utils v0.0.0
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/grpc v1.33.1
)

replace github.com/MaibornWolff/maDocK8s/core/types => ../../types

replace github.com/MaibornWolff/maDocK8s/core/utils => ../../utils
