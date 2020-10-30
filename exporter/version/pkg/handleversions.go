package pkg

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/exporter/version/adapter"
	"github.com/pkg/errors"
)

func (s *service) Notify(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	s.logger.Info("Got Notify message: ", instance)

	return s.Execute(ctx, instance)
}

func (s *service) NotifyDelete(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	s.logger.Infof("Got NotifyDelete message: %v", instance)

	return s.Execute(ctx, instance)
}

func (s *service) Execute(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	namespace := instance.GetNamespace()

	deployments, err := s.getVersions(namespace)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	markdown, err := s.markdownFor(namespace, deployments)
	if err != nil {
		err = errors.Wrapf(err, "cannot create markdown")
		s.logger.Error(err)
		return nil, err
	}

	response, err := s.storage.StoreMarkdown(ctx, markdown)
	if err != nil {
		err = errors.Wrapf(err, "cannot store markdown")
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Infof("storage stored file: %s", response.Filename)

	return new(notifier.NotifierResult), nil
}

func (s *service) markdownFor(namespace string, deployments []adapter.Deployment) (*protocol.Markdown, error) {
	content, err := s.generateContent(namespace, deployments)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate content ")
	}

	markdown := &protocol.Markdown{
		Content:  content,
		Name:     fmt.Sprintf("ns=%s", namespace),
		Source:   namespace,
		Exporter: "Versions",
		Tags:     []string{"version", "containers", "images"},
	}

	return markdown, nil
}

func (s *service) generateContent(namespace string, deployments []adapter.Deployment) (string, error) {
	var contentResult bytes.Buffer

	t, err := template.ParseFiles(s.templateFile)
	if err != nil {
		return "", errors.Wrap(err, "cannot load template file declaration")
	}

	t.Execute(&contentResult, &struct {
		Namespace   string
		Deployments []adapter.Deployment
	}{Namespace: namespace, Deployments: deployments})

	return contentResult.String(), nil
}
