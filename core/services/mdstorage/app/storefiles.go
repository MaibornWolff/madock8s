package app

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/types/services/mdstorage"

	"github.com/MaibornWolff/maDocK8s/core/utils/errors"
)

var baseFolder string = "docs"

func (s *Service) StoreMarkdown(ctx context.Context, markdown *protocol.Markdown) (*mdstorage.StoreResult, error) {
	s.logger.Infof("STORE markdown: %s for %s", markdown.Exporter, markdown.Name)
	if err := validate(markdown); err != nil {
		return nil, err
	}

	content, filename := handle(markdown)
	size, err := s.fileOutputProvider.WriteFile(content, filename)

	return &mdstorage.StoreResult{Size: int32(size), Filename: filename}, err
}

func handle(markdown *protocol.Markdown) (string, string) {
	content := prepareFileContent(markdown)
	filename := prepareFileName(markdown)

	return content, path.Join(baseFolder, filename)
}

func prepareFileContent(markdown *protocol.Markdown) string {
	time := time.Now().Format(time.StampMilli)
	tags := strings.Join(markdown.Tags, ", ")
	content := fmt.Sprintf("Generated on: %s\n%s\nTags: %s", time, markdown.Content, tags)
	return content
}

func replaceSpecialCharacters(input string) string {
	chars := []string{"]", "^", "\\\\", "[", "(", ")", ":", "/", "-"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	return re.ReplaceAllString(input, "")
}

func prepareFileName(markdown *protocol.Markdown) string {
	service := strings.Title(markdown.Name)
	exporter := strings.Title(markdown.Exporter)
	name := fmt.Sprintf("%s_%s.md", service, exporter)
	return replaceSpecialCharacters(name)
}

func validate(markdown *protocol.Markdown) error {
	if markdown == nil {
		return &errors.ErrInvalidInput{Argument: "Markdown", Method: "MdStorage.StoreMarkdown"}
	}

	if markdown.Content == "" {
		return &errors.ErrInvalidInput{Argument: "Content", Method: "MdStorage.StoreMarkdown"}
	}

	if markdown.Source == "" {
		return &errors.ErrInvalidInput{Argument: "Source", Method: "MdStorage.StoreMarkdown"}
	}

	if markdown.Name == "" {
		return &errors.ErrInvalidInput{Argument: "Name", Method: "MdStorage.StoreMarkdown"}
	}

	if markdown.Exporter == "" {
		return &errors.ErrInvalidInput{Argument: "Exporter", Method: "MdStorage.StoreMarkdown"}
	}

	return nil
}
