package pkg

import (
	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/sirupsen/logrus"
)

type service struct {
	templateFile string
	storage      mdstorage.MdstorageServer
	endpoint     func(address string) (string, error)
	logger       *logrus.Entry
	name         string
}

func NewService(name string, storage mdstorage.MdstorageServer, endpoint func(address string) (string, error), logger *logrus.Entry) notifier.NotifierServer {
	template := "/etc/prometheus/static/template.md"

	return &service{
		templateFile: template,
		storage:      storage,
		endpoint:     endpoint,
		logger:       logger,
		name:         name,
	}
}
