module github.com/MaibornWolff/maDocK8s/exporter/github

go 1.15

require (
	github.com/MaibornWolff/maDocK8s/core/types v0.0.0
	github.com/MaibornWolff/maDocK8s/core/utils v0.0.0-00010101000000-000000000000
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200806125547-5acd03effb82 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.33.1
)

replace github.com/MaibornWolff/maDocK8s/core/types => ../../core/types

replace github.com/MaibornWolff/maDocK8s/core/utils => ../../core/utils
