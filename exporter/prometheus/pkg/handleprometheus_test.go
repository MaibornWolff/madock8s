package pkg

import (
	"context"
	"strings"
	"testing"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/sirupsen/logrus"

	. "github.com/onsi/gomega"
)

func TestShouldStoreEveryGeneratedMarkdownForAddress(t *testing.T) {

	metricblocktext := `# HELP jvm_memory_bytes_used Used bytes of a given JVM memory area.
	# TYPE jvm_memory_bytes_used gauge
	jvm_memory_bytes_used{area="heap",} 1.13428264E8`

	expContent := `# Metrics from service __serviceName__ at __url__

**jvm_memory_bytes_used as gauge**
**Used bytes of a given JVM memory area.**
`
	expContent1 := strings.Replace(expContent, "__url__", "http://service1/", -1)
	expContent1 = strings.Replace(expContent1, "__serviceName__", "MetricsService1", -1)
	expItem1 := protocol.Markdown{
		Content:  expContent1,
		Name:     "MetricsService1",
		Source:   "http://service1/",
		Exporter: "Prometheus",
		Tags:     []string{"metrics", "prometheus", "monitoring"},
	}

	expContent2 := strings.Replace(expContent, "__url__", "http://service2/", -1)
	expContent2 = strings.Replace(expContent2, "__serviceName__", "MetricsService2", -1)
	expItem2 := protocol.Markdown{
		Content:  expContent2,
		Name:     "MetricsService2",
		Source:   "http://service2/",
		Exporter: "Prometheus",
		Tags:     []string{"metrics", "prometheus", "monitoring"},
	}

	endPointMock := func(string) (string, error) { return metricblocktext, nil }
	documents := []protocol.Markdown{}
	Ω := NewGomegaWithT(t)
	instance1 := &notifier.Instance{
		Name:    "MetricsService1",
		Address: "http://service1/",
	}
	instance2 := &notifier.Instance{
		Name:    "MetricsService2",
		Address: "http://service2/",
	}

	mockedMdstorageServer := &mdstorage.MdstorageServerMock{
		StoreMarkdownFunc: func(ctx context.Context, md *protocol.Markdown) (*mdstorage.StoreResult, error) {
			documents = append(documents, *md)
			return nil, nil
		},
	}

	name := "prometheus"
	log := logrus.WithField("exporter", name)
	sut := NewService(name, mockedMdstorageServer, endPointMock, log).(*service)
	sut.templateFile = "../static/template_test_mock.md"
	sut.Notify(nil, instance1)
	sut.Notify(nil, instance2)

	Ω.Expect(len(documents)).To(Equal(2))
	Ω.Expect(documents).Should(ConsistOf(expItem1, expItem2))
}
