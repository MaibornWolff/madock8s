module github.com/MaibornWolff/maDocK8s/exporter/version

go 1.15

require (
	github.com/MaibornWolff/maDocK8s/core/types v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/grpc v1.32.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/MaibornWolff/maDocK8s/core/types => ../../core/types
