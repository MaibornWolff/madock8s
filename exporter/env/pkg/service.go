package pkg

import (
	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/MaibornWolff/maDocK8s/exporter/env/adapter"
	"github.com/sirupsen/logrus"
)

type service struct {
	templateFile string
	storage      mdstorage.MdstorageServer
	environment  func(address string, namespace string) (adapter.ContainerVarsMap, error)
	logger       *logrus.Entry
	name         string
}

func NewService(name string, storage mdstorage.MdstorageServer, environment func(string, string) (adapter.ContainerVarsMap, error), logger *logrus.Entry) notifier.NotifierServer {
	template := "/etc/env-exporter/static/template.md"

	return &service{
		templateFile: template,
		storage:      storage,
		environment:  environment,
		logger:       logger,
		name:         name,
	}
}
