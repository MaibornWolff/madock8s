package adapter

import (
	"strings"

	"github.com/MaibornWolff/maDocK8s/exporter/version/utils"
	"github.com/pkg/errors"
	"k8s.io/api/apps/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	Name       string
	Containers []Container
}
type Container struct {
	Name    string
	Image   string
	Version string
}

func GetVersions(namespace string) ([]Deployment, error) {
	kubeclient := GetClient()

	deploymentList, err := kubeclient.AppsV1beta1().Deployments(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed GET for deployments")
	}

	if utils.INCLUDE_LABELED_ONLY() == true {
		deploymentList = filterDeploymentList(deploymentList)
	}

	deployments := FlattenDeployments(deploymentList)

	return deployments, nil
}

func filterDeploymentList(deploymentList *v1beta1.DeploymentList) *v1beta1.DeploymentList {
	filteredList := v1beta1.DeploymentList{
		ListMeta: deploymentList.ListMeta,
		TypeMeta: deploymentList.TypeMeta,
		Items:    []v1beta1.Deployment{},
	}
	for _, deploy := range deploymentList.Items {
		if deploy.Annotations["madock8s.exporter/versionExporter"] == "true" {
			filteredList.Items = append(filteredList.Items, deploy)
		}
	}
	return &filteredList
}

func FlattenDeployments(deploymentList *v1beta1.DeploymentList) []Deployment {
	deployments := []Deployment{}
	for _, deploy := range deploymentList.Items {
		containers := []Container{}
		for _, c := range deploy.Spec.Template.Spec.Containers {
			imageSplit := strings.Split(c.Image, ":")

			image := imageSplit[0]
			version := imageSplit[1]
			container := Container{
				Name:    c.Name,
				Image:   image,
				Version: version,
			}
			containers = append(containers, container)
		}
		deployments = append(deployments, Deployment{
			Name:       deploy.Name,
			Containers: containers,
		})
	}

	return deployments
}
