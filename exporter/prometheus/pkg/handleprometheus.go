package pkg

import (
	"bytes"
	"context"
	"text/template"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/utils/envutils"

	"github.com/pkg/errors"
)

func (s *service) Notify(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	s.logger.Info("Got Notify message :", instance)
	address := instance.GetAddress()
	serviceName := instance.GetName()
	markdown, err := s.markdownForAddress(address, serviceName)
	if err != nil {
		err = errors.Wrapf(err, "cannot generate markdown for %v", address)
		s.logger.Error(err)
		return nil, err
	}

	response, err := s.storage.StoreMarkdown(ctx, markdown)
	if err != nil {
		err = errors.Wrapf(err, "cannot store markdown for %v", address)
		s.logger.Error(err)
		return nil, err
	}
	s.logger.Infof("storage stored file: %s", response.Filename)

	return new(notifier.NotifierResult), nil
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

func (s *service) newMarkdown(content string, name string, address string) *protocol.Markdown {
	return &protocol.Markdown{
		Content:  content,
		Name:     name,
		Source:   address,
		Exporter: s.name,
		Tags:     []string{"metrics", "prometheus", "monitoring"},
	}
}

func (s *service) markdownForAddress(address string, serviceName string) (*protocol.Markdown, error) {
	body, err := s.endpoint(address)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get request body for markdown")
	}
	parsedResult := ParsePrometheusResult(body)
	content, err := s.generateContent(serviceName, address, parsedResult)
	if err != nil {
		return nil, errors.Wrap(err, "generation of content failed")
	}

	markdown := s.newMarkdown(content, serviceName, address)

	return markdown, nil
}

func (s *service) generateContent(name string, address string, metrics []MetricsMetaBlock) (string, error) {
	var contentResult bytes.Buffer

	t, err := template.ParseFiles(s.templateFile)
	if err != nil {
		return "", errors.Wrap(err, "cannot load template file declaration")
	}

	t.Execute(&contentResult, &struct {
		Name    string
		Address string
		Blocks  []MetricsMetaBlock
	}{Name: name, Address: address, Blocks: metrics})

	return contentResult.String(), nil
}
