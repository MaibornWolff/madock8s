package pkg

import (
	"bytes"
	"context"
	"html/template"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/utils/envutils"
	"github.com/MaibornWolff/maDocK8s/exporter/env/adapter"
	"github.com/pkg/errors"
)

type uInstance struct {
	name           string
	address        string
	deploymentName string
	namespace      string
}

func (s *service) Notify(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	s.logger.Infof("Got Notify message: %v", instance)

	uInstance := unwrapInstance(instance)

	markdown, err := s.markdownForAddress(uInstance)
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
		Tags:     []string{"environment", "env"},
	}
}

func unwrapInstance(instance *notifier.Instance) uInstance {
	namespace := instance.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	config := notifier.ConvertNotifierMapToStringMap(instance.Configuration)

	return uInstance{
		name:           instance.GetName(),
		namespace:      namespace,
		address:        instance.GetAddress(),
		deploymentName: config["value"],
	}
}

func (s *service) markdownForAddress(instance uInstance) (*protocol.Markdown, error) {
	envVars, err := s.environment(instance.deploymentName, instance.namespace)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get request body for markdown")
	}

	content, err := s.generateContent(instance, envVars)
	markdown := s.newMarkdown(content, instance.name, instance.address)
	return markdown, nil
}

func (s *service) generateContent(instance uInstance, containerVarsMap adapter.ContainerVarsMap) (string, error) {
	var contentResult bytes.Buffer

	t, err := template.ParseFiles(s.templateFile)
	if err != nil {
		return "", errors.Wrap(err, "cannot load template file declaration")
	}

	t.Execute(&contentResult, &struct {
		Name             string
		Address          string
		DeploymentName   string
		ContainerVarsMap adapter.ContainerVarsMap
	}{Name: instance.name, Address: instance.address, DeploymentName: instance.deploymentName, ContainerVarsMap: containerVarsMap})

	return contentResult.String(), nil
}
