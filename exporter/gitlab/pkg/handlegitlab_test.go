package pkg

import (
	"context"
	"strings"
	"testing"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/MaibornWolff/maDocK8s/exporter/gitlab/adapter"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestShouldStoreGeneratedMarkdownForGitLabUrl(t *testing.T) {
	file1Content := `##### Example Readme File
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`

	file2Content := `##### Example Readme File 2
Lorem ipsum dolor sit amet, consectetur adipiscing elit.`
	expContent := `# Markdown files for service __name__ at __url__
## These are all markdown files found in repository of the service

Readme.md

##### Example Readme File
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

folder/Nested.md

##### Example Readme File 2
Lorem ipsum dolor sit amet, consectetur adipiscing elit.

`

	expContent1 := strings.Replace(expContent, "__url__", "http://service.com/", -1)
	expContent1 = strings.Replace(expContent1, "__name__", "ExampleService", -1)

	expItem := protocol.Markdown{
		Content:  expContent1,
		Name:     "ExampleService",
		Source:   "http://service.com/",
		Exporter: "GitLab",
		Tags:     []string{"readme", "gitlab", "information"},
	}

	mockedGetMdFiles := func(map[string]string) ([]adapter.GitlabFile, error) {
		return []adapter.GitlabFile{
			{
				Name:    "Readme.md",
				Path:    "Readme.md",
				Content: file1Content,
			}, {
				Name:    "Nested.md",
				Path:    "folder/Nested.md",
				Content: file2Content,
			},
		}, nil
	}

	document := protocol.Markdown{}
	Ω := NewWithT(t)

	config := map[string]string{
		"baseurl": "https://git.example.com",
		"path":    "service",
		"ref":     "master",
	}

	instance := &notifier.Instance{
		Name:          "ExampleService",
		Address:       "http://service.com/",
		Configuration: notifier.ConvertStringMapToNotifierMap(config),
	}

	mockedMdStorageServer := &mdstorage.MdstorageServerMock{
		StoreMarkdownFunc: func(ctx context.Context, md *protocol.Markdown) (*mdstorage.StoreResult, error) {
			document = *md
			return nil, nil
		},
	}

	name := "gitlab"
	log := logrus.WithField("exporter", name)
	sut := NewService(name, mockedMdStorageServer, mockedGetMdFiles, log).(*service)
	sut.templateFile = "../static/template-test.md"
	sut.Notify(nil, instance)

	Ω.Expect(&document).Should(Equal(&expItem))
}
