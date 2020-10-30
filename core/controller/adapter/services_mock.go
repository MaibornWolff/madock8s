package adapter

import (
	"errors"

	api_v1 "k8s.io/api/core/v1"
)

type ServicesMock struct {
	ServiceList *api_v1.ServiceList
	ID          string
}

func (mock *ServicesMock) GetServices(madock8sID string, namespace string) (*api_v1.ServiceList, error) {
	if madock8sID == mock.ID {
		return mock.ServiceList, nil
	}
	return nil, errors.New("not found")
}
