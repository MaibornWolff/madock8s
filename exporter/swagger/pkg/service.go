package pkg

import (
	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/MaibornWolff/maDocK8s/exporter/swagger/adapter"
	"github.com/sirupsen/logrus"
)

type service struct {
	templateFile   string
	storage        mdstorage.MdstorageServer
	getSwaggerSpec func(string) (map[string][]adapter.Endpoint, error)
	logger         *logrus.Entry
	name           string
}

func NewService(name string, storage mdstorage.MdstorageServer, getSwaggerSpec func(string) (map[string][]adapter.Endpoint, error), logger *logrus.Entry) notifier.NotifierServer {
	template := "/etc/swagger/static/template.md"

	return &service{
		templateFile:   template,
		storage:        storage,
		getSwaggerSpec: getSwaggerSpec,
		logger:         logger,
		name:           name,
	}
}
