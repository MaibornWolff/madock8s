package adapter

import (
	"context"
	"errors"

	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"google.golang.org/grpc/grpclog"

	"google.golang.org/grpc"
)

type MdStorageProvider struct {
	Address string
}

func (s *MdStorageProvider) StoreMarkdown(ctx context.Context, markdown *protocol.Markdown) (*mdstorage.StoreResult, error) {
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := mdstorage.NewMdstorageClient(conn)

	response, err := client.StoreMarkdown(ctx, markdown)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	return response, err
}

func (s *MdStorageProvider) UpdateMarkdown(ctx context.Context, markdown *protocol.Markdown) (*mdstorage.StoreResult, error) {
	return nil, errors.New("Not allowed for this exporter")
}

func (s *MdStorageProvider) DeleteMarkdown(ctx context.Context, markdown *protocol.Markdown) (*mdstorage.RemoveResult, error) {
	return nil, errors.New("Not allowed for this exporter")
}
