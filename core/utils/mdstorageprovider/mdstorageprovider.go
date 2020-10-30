package mdstorageprovider

import (
	"context"

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
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := mdstorage.NewMdstorageClient(conn)

	response, err := client.UpdateMarkdown(ctx, markdown)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	return response, err
}

func (s *MdStorageProvider) DeleteMarkdown(ctx context.Context, markdown *protocol.Markdown) (*mdstorage.RemoveResult, error) {
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := mdstorage.NewMdstorageClient(conn)

	response, err := client.DeleteMarkdown(ctx, markdown)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
	}
	return response, err
}
