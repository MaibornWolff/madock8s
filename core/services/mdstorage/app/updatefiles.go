package app

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"
	"github.com/pkg/errors"
)

func (s *Service) UpdateMarkdown(ctx context.Context, markdown *protocol.Markdown) (*mdstorage.StoreResult, error) {
	s.logger.Infof("UPDATE markdown: %s for %s", markdown.Exporter, markdown.Name)

	filename := prepareFileName(markdown)
	filename = path.Join(baseFolder, filename)

	storedContent, err := s.fileOutputProvider.ReadFile(filename)
	if err != nil {
		err = errors.Wrapf(err, "cannot read stored file")
		s.logger.Error(err)
		return nil, err
	}

	newContent := prepareContent(storedContent)
	size, err := s.fileOutputProvider.WriteFile(newContent, filename)
	return &mdstorage.StoreResult{Size: int32(size), Filename: filename}, err
}

func prepareContent(storedContent string) string {
	time := time.Now().Format(time.StampMilli)
	return fmt.Sprintf("**Deployments of the service were deleted on %s**\n\n%s", time, storedContent)
}

func (s *Service) DeleteMarkdown(ctx context.Context, markdown *protocol.Markdown) (*mdstorage.RemoveResult, error) {
	s.logger.Infof("DELETE markdown: %s for %s", markdown.Exporter, markdown.Name)
	filename := prepareFileName(markdown)
	filename = path.Join(baseFolder, filename)
	err := s.fileOutputProvider.RemoveFile(filename)
	return &mdstorage.RemoveResult{}, err
}
