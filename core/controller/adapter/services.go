package adapter

import (
	"github.com/MaibornWolff/maDocK8s/core/controller/utils"
	"github.com/pkg/errors"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func GetServices(madock8sID string, namespace string) (*api_v1.ServiceList, error) {
	label := labels.SelectorFromSet(labels.Set(map[string]string{"madock8s": madock8sID}))
	options := meta_v1.ListOptions{
		LabelSelector: label.String(),
	}

	var kubeClient kubernetes.Interface
	kubeClient = utils.GetClient()
	services, err := kubeClient.CoreV1().Services(namespace).List(options)
	if err != nil {
		err = errors.Wrapf(err, "Error from api core v1 services: %v", err)
		return nil, err
	}

	return services, nil
}
