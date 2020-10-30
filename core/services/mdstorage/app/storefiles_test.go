package app

import (
	"testing"

	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/sirupsen/logrus"

	"github.com/MaibornWolff/maDocK8s/core/services/mdstorage/adapter"

	. "github.com/onsi/gomega"
)

var success = new(adapter.FileOutputProviderSuccessMock)
var writeerror = new(adapter.FileOutputProviderErrorMock)

var markdownTable = []struct {
	input    *protocol.Markdown
	output   *mdstorage.StoreResult
	errs     bool
	provider adapter.FileOutputProvider
}{
	{nil, nil, true, success},
	{&protocol.Markdown{Name: ""}, nil, true, success},
	{&protocol.Markdown{
		Name:     "myservice",
		Content:  "*metrics*",
		Source:   "Prometheus",
		Exporter: "PrometheusExporter",
	}, &mdstorage.StoreResult{Size: 100}, false, success},
	{&protocol.Markdown{
		Name:     "myservice",
		Content:  "*metrics*",
		Source:   "Prometheus",
		Exporter: "PrometheusExporter",
	}, &mdstorage.StoreResult{Size: 0}, true, writeerror},
}

func TestStoreMarkdown(t *testing.T) {
	Ω := NewGomegaWithT(t)
	for _, params := range markdownTable {
		log := logrus.WithField("service", "mdstorage")
		s := NewService(params.provider, log)

		output, err := s.StoreMarkdown(nil, params.input)
		if params.errs {
			Ω.Expect(err).To(HaveOccurred())
			Ω.Expect(output).To(Equal(params.output))
		} else {
			Ω.Expect(err).ToNot(HaveOccurred())
			Ω.Expect(output).To(Equal(params.output))
		}
	}
}

func TestHandleMarkdown(t *testing.T) {
	Ω := NewGomegaWithT(t)

	markdown := &protocol.Markdown{
		Name:     "myservice",
		Content:  "*metrics*",
		Source:   "Prometheus",
		Tags:     []string{"tag1", "tag2"},
		Exporter: "PrometheusExporter",
	}

	content, filename := handle(markdown)
	Ω.Expect(content).To(Equal("*metrics*\nTags:tag1,tag2"))
	Ω.Expect(filename).To(Equal("docs/PrometheusExporter_Prometheus_myservice.md"))

}

func TestShouldEscapeAddressForFileNameMarkdown(t *testing.T) {
	Ω := NewGomegaWithT(t)

	markdown := &protocol.Markdown{
		Name:     "myservice",
		Content:  "*metrics*",
		Source:   "http://thisisaaddress:8080/",
		Tags:     []string{"tag1", "tag2"},
		Exporter: "PrometheusExporter",
	}

	content, filename := handle(markdown)
	Ω.Expect(content).To(Equal("*metrics*\nTags:tag1,tag2"))
	Ω.Expect(filename).To(Equal("docs/PrometheusExporter_httpthisisaaddress8080_myservice.md"))

}
