package pkg

import (
	"context"
	"strings"
	"testing"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/MaibornWolff/maDocK8s/exporter/env/adapter"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestShouldStoreGeneratedMarkdownForGitLabUrl(t *testing.T) {
	expContent := `# ENV for service __name__ at __url__
## This is a list of environment variables for deployment __deploy__<br/>
| Key | SourceType/Name.Key | Value |
|---|:-----------:|-----:|
| __Container *` + "`" + "container-name-1" + "`" + `*__|||<BR />
| ENV_VAR_1 |  | Env variable 1 |<BR />
| REFERENCED_VAR | configmap/example-config-map.var | Referenced value |<BR />
| __Container *` + "`" + "container-name-2" + "`" + `*__|||<BR />
| ENV_VAR_2 |  | Env variable 2 |<BR />

`

	serviceName := "ExampleService"
	deploymentName := "example-service"
	url := "http://service.com/"

	expContent = strings.Replace(expContent, "__name__", serviceName, -1)
	expContent = strings.Replace(expContent, "__url__", url, -1)
	expContent = strings.Replace(expContent, "__deploy__", deploymentName, -1)

	expItem := protocol.Markdown{
		Content:  expContent,
		Name:     serviceName,
		Source:   url,
		Exporter: "ENV",
		Tags:     []string{"environment", "env"},
	}

	mockedEnvironment := adapter.ContainerVarsMap{}
	mockedEnvironment["container-name-1"] = []adapter.EnvVar{
		{
			Key:   "ENV_VAR_1",
			Value: "Env variable 1",
		},
		{
			Key:           "REFERENCED_VAR",
			ExtSourceType: "configmap",
			ExtSourceName: "example-config-map",
			ExtSourceKey:  "var",
			Value:         "Referenced value",
		},
	}
	mockedEnvironment["container-name-2"] = []adapter.EnvVar{
		{
			Key:   "ENV_VAR_2",
			Value: "Env variable 2",
		},
	}

	mockedGetEnvironmentVars := func(string, string) (adapter.ContainerVarsMap, error) { return mockedEnvironment, nil }

	document := protocol.Markdown{}
	Ω := NewWithT(t)

	config := notifier.ConvertStringMapToNotifierMap(map[string]string{"value": deploymentName})
	instance := &notifier.Instance{
		Name:          "ExampleService",
		Address:       "http://service.com/",
		Configuration: config,
	}

	mockedMdStorageServer := &mdstorage.MdstorageServerMock{
		StoreMarkdownFunc: func(ctx context.Context, md *protocol.Markdown) (*mdstorage.StoreResult, error) {
			document = *md
			return nil, nil
		},
	}

	name := "environment"
	log := logrus.WithField("exporter", name)
	sut := NewService(name, mockedMdStorageServer, mockedGetEnvironmentVars, log).(*service)
	sut.templateFile = "../static/template-test.md"
	sut.Notify(nil, instance)

	Ω.Expect(&document).Should(Equal(&expItem))
}
