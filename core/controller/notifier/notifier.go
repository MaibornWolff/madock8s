package notifier

import (
	"context"

	notifier "github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type Exporter struct {
	Name          string
	Configuration map[string]string
	ClusterIP     string
}

func NotifyExporter(exporter Exporter, address string, serviceName string, namespace string) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(exporter.ClusterIP, opts...)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := notifier.NewNotifierClient(conn)

	instance := &notifier.Instance{
		Name:          serviceName,
		Namespace:     namespace,
		Address:       address,
		Configuration: notifier.ConvertStringMapToNotifierMap(exporter.Configuration),
	}

	_, err = client.Notify(context.Background(), instance)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
}

func NotifyDeleteExporter(exporter Exporter, serviceName string, namespace string) {
	conn, err := grpc.Dial(exporter.ClusterIP, grpc.WithInsecure())
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := notifier.NewNotifierClient(conn)
	instance := &notifier.Instance{
		Name:          serviceName,
		Namespace:     namespace,
		Configuration: notifier.ConvertStringMapToNotifierMap(exporter.Configuration),
	}
	_, err = client.NotifyDelete(context.Background(), instance)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
}
