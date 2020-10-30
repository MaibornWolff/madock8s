package pkg

import (
	"context"
	"strings"
	"testing"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/MaibornWolff/maDocK8s/exporter/swagger/adapter"
	"github.com/go-openapi/loads"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestShouldStoreGeneratedMarkdown(t *testing.T) {
	expContent := `# Swagger-documented endpoints for service __name__ at __address__

- Swagger.json location: __json_url__
- Swagger home page (internal access only): __base_url__

## List of available endpoints:

| Path | Method | Summary |
| :---: | :----: | ------: |
| /greetings | POST | Greetings returns a greeting to the developer. |
| /metrics | GET | Metrics returns a list of metrics available for this service. |

`

	expContent = strings.Replace(expContent, "__name__", "ExampleService", -1)
	expContent = strings.Replace(expContent, "__address__", "192.168.0.1:80", -1)
	expContent = strings.Replace(expContent, "__json_url__", "http://192.168.0.1:81/docs/swagger.json", -1)
	expContent = strings.Replace(expContent, "__base_url__", "http://192.168.0.1:81/docs", -1)

	expItem := protocol.Markdown{
		Content:  expContent,
		Name:     "ExampleService",
		Source:   "192.168.0.1:80",
		Exporter: "Swagger",
		Tags:     []string{"API", "swagger", "docs"},
	}

	mockedGetSwaggerSpec := func(string) (map[string][]adapter.Endpoint, error) {
		doc, err := loads.JSONSpec("../static/swagger-test.json")
		if err != nil {
			return nil, err
		}
		return adapter.ExtractEndpoints(doc), nil
	}

	document := protocol.Markdown{}
	Ω := NewWithT(t)

	config := map[string]string{
		"baseurl": "/docs",
		"jsonurl": "",
		"json":    "/swagger.json",
		"port":    "81",
	}

	instance := &notifier.Instance{
		Name:          "ExampleService",
		Address:       "192.168.0.1:80",
		Configuration: notifier.ConvertStringMapToNotifierMap(config),
	}

	mockedMdStorageServer := &mdstorage.MdstorageServerMock{
		StoreMarkdownFunc: func(ctx context.Context, md *protocol.Markdown) (*mdstorage.StoreResult, error) {
			document = *md
			return nil, nil
		},
	}
	name := "swagger"
	log := logrus.WithField("exporter", name)
	sut := NewService(name, mockedMdStorageServer, mockedGetSwaggerSpec, log).(*service)
	sut.templateFile = "../static/template-test.md"
	sut.Notify(nil, instance)

	Ω.Expect(&document).Should(Equal(&expItem))
}
