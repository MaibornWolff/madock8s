package pkg

import (
	"context"
	"strings"
	"testing"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/MaibornWolff/maDocK8s/exporter/version/adapter"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var expDeployment1 = adapter.Deployment{
	Name: "example-deployment",
	Containers: []adapter.Container{
		adapter.Container{
			Name:    "example-container-1",
			Image:   "example.com/image-1",
			Version: "v1",
		},
		adapter.Container{
			Name:    "example-container-2",
			Image:   "example.com/image-2",
			Version: "1.1.1",
		},
	},
}

var expDeployment2 = adapter.Deployment{
	Name: "other-deployment",
	Containers: []adapter.Container{
		adapter.Container{
			Name:    "other-container",
			Image:   "example.com/other-image",
			Version: "3.0.1",
		},
	},
}

func TestShouldStoreAllImageVersionsForNamespace(t *testing.T) {
	namespace := "staging"
	expContent := `# Image versions for deployments from namespace __namespace__

| Deployment | ContainerName | ImageName | ImageVersion |
| :--------- | :-----------: | :-------: | -----------: |
| example-deployment | example-container-1 | example.com/image-1 | v1 |
| example-deployment | example-container-2 | example.com/image-2 | 1.1.1 |
| other-deployment | other-container | example.com/other-image | 3.0.1 |

`
	expContent = strings.Replace(expContent, "__namespace__", namespace, -1)

	expItem := protocol.Markdown{
		Content:  expContent,
		Name:     namespace,
		Source:   namespace,
		Exporter: "Versions",
		Tags:     []string{"version", "containers", "images"},
	}

	mockedGetVersions := func(string) ([]adapter.Deployment, error) {
		return adapter.FlattenDeployments(fetchDeploymentsMock()), nil
	}

	document := protocol.Markdown{}

	instance := &notifier.Instance{
		Name:      "ExampleService",
		Address:   "http://service.com/",
		Namespace: namespace,
	}

	mockedMdStorageServer := &mdstorage.MdstorageServerMock{
		StoreMarkdownFunc: func(ctx context.Context, md *protocol.Markdown) (*mdstorage.StoreResult, error) {
			document = *md
			return nil, nil
		},
	}

	name := "version"
	log := logrus.WithField("exporter", name)
	sut := NewService(name, mockedMdStorageServer, mockedGetVersions, log).(*service)
	sut.templateFile = "../static/template-test.md"
	sut.Notify(nil, instance)

	Ω := NewWithT(t)
	Ω.Expect(&document).Should(Equal(&expItem))
}

func TestShouldExtractVersionCorrectly(t *testing.T) {
	result := adapter.FlattenDeployments(fetchDeploymentsMock())

	Ω := NewWithT(t)
	Ω.Expect(len(result)).To(Equal(2))
	Ω.Expect(result).Should(ConsistOf(expDeployment1, expDeployment2))
}

func fetchDeploymentsMock() *v1beta1.DeploymentList {
	deployments := &v1beta1.DeploymentList{
		Items: []v1beta1.Deployment{
			v1beta1.Deployment{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "example-deployment",
					Annotations: map[string]string{
						"madock8s.exporter/versionExporter": "true",
					},
				},
				Spec: v1beta1.DeploymentSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Image: "example.com/image-1:v1",
									Name:  "example-container-1",
								},
								v1.Container{
									Image: "example.com/image-2:1.1.1",
									Name:  "example-container-2",
								},
							},
						},
					},
				},
			},
			v1beta1.Deployment{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "other-deployment",
					Annotations: map[string]string{
						"madock8s.exporter/versionExporter": "true",
					},
				},
				Spec: v1beta1.DeploymentSpec{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Image: "example.com/other-image:3.0.1",
									Name:  "other-container",
								},
							},
						},
					},
				},
			},
		},
	}
	return deployments
}
