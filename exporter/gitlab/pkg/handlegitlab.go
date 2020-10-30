package pkg

import (
	"bytes"
	"context"
	"html/template"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/utils/envutils"
	"github.com/MaibornWolff/maDocK8s/exporter/gitlab/adapter"
	"github.com/pkg/errors"
)

type uInstance struct {
	name    string
	address string
	config  map[string]string
	mdFiles []adapter.GitlabFile
}

func (s *service) NotifyDelete(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	strategy := envutils.GetDeletionStrategyFromEnv()
	s.logger.Infof("Got NotifyDelete %s message: %v", strategy, instance)
	var err error
	switch strategy {
	case "IGNORE":
		return new(notifier.NotifierResult), nil
	case "UPDATE":
		markdown := s.newMarkdown("", instance.Name, instance.Address)
		_, err = s.storage.UpdateMarkdown(ctx, markdown)
	case "DELETE":
		markdown := s.newMarkdown("", instance.Name, instance.Address)
		_, err = s.storage.DeleteMarkdown(ctx, markdown)
	}
	if err != nil {
		err = errors.Wrapf(err, "cannot %s markdown for %v", strategy, instance.Name)
		s.logger.Error(err)
	}
	return new(notifier.NotifierResult), err
}

func (s *service) Notify(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	s.logger.Infof("Got Notify message: %v", instance)

	uInstance := uInstance{
		name:    instance.GetName(),
		address: instance.GetAddress(),
		config:  notifier.ConvertNotifierMapToStringMap(instance.GetConfiguration()),
	}

	mdFiles, err := s.mdFiles(uInstance.config)
	if err != nil {
		err = errors.Wrapf(err, "cannot fetch files for %v at %v", uInstance.name, uInstance.address)
		s.logger.Error(err)
		return nil, err
	}

	uInstance.mdFiles = mdFiles

	markdown, err := s.generateMarkdownFor(uInstance)
	if err != nil {
		err = errors.Wrapf(err, "cannot generate markdown for %v at %v", uInstance.name, uInstance.address)
		s.logger.Error(err)
		return nil, err
	}

	response, err := s.storage.StoreMarkdown(ctx, markdown)
	if err != nil {
		err = errors.Wrapf(err, "cannot store markdown for %v at %v", uInstance.name, uInstance.address)
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Infof("storage stored file: %s", response.Filename)

	return new(notifier.NotifierResult), nil
}

func (s *service) newMarkdown(content string, name string, address string) *protocol.Markdown {
	return &protocol.Markdown{
		Content:  content,
		Name:     name,
		Source:   address,
		Exporter: s.name,
		Tags:     []string{"readme", "gitlab", "information"},
	}
}

func (s *service) generateMarkdownFor(instance uInstance) (*protocol.Markdown, error) {
	content, _ := s.generateContent(instance)
	markdown := s.newMarkdown(content, instance.name, instance.address)
	return markdown, nil
}

func (s *service) generateContent(instance uInstance) (string, error) {
	var contentResult bytes.Buffer

	t, err := template.ParseFiles(s.templateFile)
	if err != nil {
		return "", errors.Wrap(err, "cannot load template file declaration")
	}

	err = t.Execute(&contentResult, &struct {
		Name    string
		Address string
		MdFiles []adapter.GitlabFile
	}{Name: instance.name, Address: instance.address, MdFiles: instance.mdFiles})
	if err != nil {
		return "", errors.Wrap(err, "cannot parse content to template")
	}

	return contentResult.String(), nil
}
