package main

import (
	"net"

	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/sirupsen/logrus"

	"github.com/MaibornWolff/maDocK8s/core/services/mdstorage/adapter"
	"github.com/MaibornWolff/maDocK8s/core/services/mdstorage/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var log *logrus.Entry

func main() {
	log = logrus.WithField("core_service", "mdstorage")

	service := app.NewService(new(adapter.OsFileOutputProvider), log)
	Start(service, ":80")
}

func Start(s mdstorage.MdstorageServer, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	mdstorage.RegisterMdstorageServer(server, s)
	reflection.Register(server)
	log.Infof("service start listen on port %v", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
