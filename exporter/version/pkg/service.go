package pkg

import (
	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/MaibornWolff/maDocK8s/exporter/version/adapter"
	"github.com/sirupsen/logrus"
)

type service struct {
	templateFile string
	storage      mdstorage.MdstorageServer
	getVersions  func(string) ([]adapter.Deployment, error)
	logger       *logrus.Entry
	name         string
}

func NewService(name string, storage mdstorage.MdstorageServer, getVersions func(string) ([]adapter.Deployment, error), logger *logrus.Entry) notifier.NotifierServer {
	template := "/etc/version/static/template.md"

	return &service{
		templateFile: template,
		storage:      storage,
		getVersions:  getVersions,
		logger:       logger,
		name:         name,
	}
}
