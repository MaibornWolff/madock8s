package pkg

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"github.com/MaibornWolff/maDocK8s/core/types/exporter/notifier"
	"github.com/MaibornWolff/maDocK8s/core/types/protocol"
	"github.com/MaibornWolff/maDocK8s/core/utils/envutils"
	"github.com/MaibornWolff/maDocK8s/exporter/swagger/adapter"
	"github.com/MaibornWolff/maDocK8s/exporter/swagger/utils"
	"github.com/pkg/errors"
)

func (s *service) newMarkdown(content string, name string, address string) *protocol.Markdown {
	return &protocol.Markdown{
		Content:  content,
		Name:     name,
		Source:   address,
		Exporter: s.name,
		Tags:     []string{"API", "swagger", "docs"},
	}
}

func (s *service) NotifyDelete(ctx context.Context, instance *notifier.Instance) (*notifier.NotifierResult, error) {
	s.logger.Infof("Got NotifyDelete message: %v", instance)
	strategy := envutils.GetDeletionStrategyFromEnv()
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

	uInstance := adapter.UInstance{
		Name:    instance.GetName(),
		Address: instance.GetAddress(),
		Config:  notifier.ConvertNotifierMapToStringMap(instance.GetConfiguration()),
	}

	var err error
	baseURL, err := buildBaseURL(uInstance.Address, uInstance.Config)
	if err != nil {
		err = errors.Wrapf(err, "failed building baseURL")
		s.logger.Error(err)
		return nil, err
	}
	uInstance.BaseURL = baseURL.String()

	jsonURL, err := getJSONURL(uInstance.BaseURL, uInstance.Config)
	if err != nil {
		err = errors.Wrapf(err, "failed building jsonURL")
		s.logger.Error(err)
		return nil, err
	}
	uInstance.JSONURL = jsonURL.String()

	uInstance.Endpoints, err = s.getSwaggerSpec(uInstance.JSONURL)
	if err != nil {
		err = errors.Wrapf(err, "cannot fetch swagger spec for %s at %s", uInstance.Name, uInstance.Address)
		s.logger.Error(err)
		return nil, err
	}

	markdown, err := s.generateMarkdownFor(uInstance)
	if err != nil {
		err = errors.Wrapf(err, "cannot generate markdown for %v at %v", uInstance.Name, uInstance.Address)
		s.logger.Error(err)
		return nil, err
	}

	response, err := s.storage.StoreMarkdown(ctx, markdown)
	if err != nil {
		err = errors.Wrapf(err, "cannot store markdown for %v at %v", uInstance.Name, uInstance.Address)
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Infof("storage stored file: %s", response.Filename)

	return new(notifier.NotifierResult), nil
}

// address is 192.168.0.1:80
func buildBaseURL(address string, config map[string]string) (*url.URL, error) {
	scheme := "http"
	host := strings.Split(address, ":")[0]
	port := strings.Split(address, ":")[1]
	docsPath := utils.GetSwaggerBaseURLFromEnv()

	if p, ok := config["port"]; ok && p != "" {
		port = p
	} else if p := utils.GetSwaggerPortFromEnv(); p != "" {
		port = p
	}

	if b, ok := config["baseurl"]; ok && b != "" {
		docsPath = b
	}
	return url.Parse(fmt.Sprintf("%s://%s:%v%s", scheme, host, port, docsPath))
}

func assembleJSONURLFromConfig(baseURL string, config map[string]string) string {
	if f, ok := config["json"]; ok && f != "" {
		return baseURL + f
	}
	return ""
}

func assembleJSONURLFromENV(baseURL string) string {
	filename := utils.GetSwaggerJSONFromEnv()
	return baseURL + filename
}

func getJSONURL(baseURL string, config map[string]string) (*url.URL, error) {
	// read config[jsonurl]
	if jsonURL, ok := config["jsonurl"]; ok && jsonURL != "" {
		return url.Parse(jsonURL)
	}

	// construct baseurl + config[docsPath] + config[json]
	if jsonURL := assembleJSONURLFromConfig(baseURL, config); jsonURL != "" {
		return url.Parse(jsonURL)
	}

	// read jsonURL from ENV
	if jsonURL := utils.GetSwaggerJsonURLFromEnv(); jsonURL != "" {
		return url.Parse(jsonURL)
	}

	// construct jsonURL from ENV
	if jsonURL := assembleJSONURLFromENV(baseURL); jsonURL != "" {
		return url.Parse(jsonURL)
	}

	return nil, errors.New("failed building JSON URL")
}

func (s *service) generateMarkdownFor(instance adapter.UInstance) (*protocol.Markdown, error) {
	content, err := s.generateContent(instance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate content")
	}

	markdown := s.newMarkdown(content, instance.Name, instance.Address)

	return markdown, nil
}

func (s *service) generateContent(instance adapter.UInstance) (string, error) {
	var content bytes.Buffer

	t, err := template.ParseFiles(s.templateFile)
	if err != nil {
		return "", errors.Wrap(err, "cannot load template file declaration")
	}

	err = t.Execute(&content, &instance)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse content to template")
	}

	return content.String(), nil
}
