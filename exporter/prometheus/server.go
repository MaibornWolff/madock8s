package main

import (
	"net"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/utils/mdstorageprovider"
	"github.com/MaibornWolff/maDocK8s/exporter/prometheus/adapter"
	prometheus "github.com/MaibornWolff/maDocK8s/exporter/prometheus/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var log *logrus.Entry

var name = "prometheus"

const port = ":81"

var (
	storageAddress = pflag.String("storage-address", "localhost:80", "Storage Address")
)

func main() {
	pflag.Parse()
	log = logrus.WithField("exporter", name)
	storageProvider := mdstorageprovider.MdStorageProvider{
		Address: *storageAddress,
	}

	prometheusService := prometheus.NewService(name, &storageProvider, adapter.GetEndpointBody, log)
	Start(prometheusService, port)
}

func Start(s notifier.NotifierServer, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	notifier.RegisterNotifierServer(server, s)
	reflection.Register(server)
	log.Infof("service start listen on port %v", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
