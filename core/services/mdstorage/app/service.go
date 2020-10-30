package app

import (
	"github.com/MaibornWolff/maDocK8s/core/services/mdstorage/adapter"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger             *logrus.Entry
	fileOutputProvider adapter.FileOutputProvider
}

func NewService(provider adapter.FileOutputProvider, logger *logrus.Entry) mdstorage.MdstorageServer {

	return &Service{
		logger:             logger,
		fileOutputProvider: provider,
	}
}
