package app

import (
	"strconv"
	"strings"

	"github.com/MaibornWolff/maDocK8s/core/controller/adapter"

	"github.com/MaibornWolff/maDocK8s/core/controller/notifier"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
)

//Annotations from Deployment:
// madock8s: "madocksIdfromDeployment"
// madock8s/nameExporter1: true or enabled or config <- NameExporter1 is the reachable servicename
// e.g.:
// madock8s/GitlabExporter: http://namemicroservicerepo

// 1) the controller is listening for updates on deployments with a Madocks id
// 2) with the annotation value, name of madock8s-deployment-id, should the service selected with an equivalent label
// e.g. madock8s:madocksIdfromDeployment. We need the reachable service url to move on.
// 3) then all madock8s/exporter should be extract from the deployment annotations
// 4) then all named exporter should be notified with the service url and the controller is done

type Madock8s struct {
	GetServices func(madock8sID string, namespace string) (*api_v1.ServiceList, error)
	Namespace   string
}

func NewMadock8s(namespace string, getServices func(madock8sID string, namespace string) (*api_v1.ServiceList, error)) *Madock8s {
	return &Madock8s{
		GetServices: getServices,
		Namespace:   namespace,
	}
}

func (madock8s *Madock8s) HandleDelete(annotations map[string]string) {
	madock8sID := getMadock8sID(annotations)
	if madock8sID == "" {
		return
	}

	exporters := madock8s.searchMadock8sExporter(annotations)
	if len(exporters) == 0 {
		return
	}
	exporters = madock8s.resolveExporters(exporters)

	madock8s.notifyDeleteExporters(exporters, madock8sID)

	adapter.CreateDauxGenerateJob()
}

func (madock8s *Madock8s) notifyDeleteExporters(exporters []notifier.Exporter, serviceName string) {
	for _, exporter := range exporters {
		notifier.NotifyDeleteExporter(exporter, serviceName, madock8s.Namespace)
	}
}

func (madock8s *Madock8s) HandleChange(annotations map[string]string) {
	madock8sID := getMadock8sID(annotations)
	if madock8sID == "" {
		return
	}

	exporters := madock8s.searchMadock8sExporter(annotations)
	if len(exporters) == 0 {
		return
	}
	exporters = madock8s.resolveExporters(exporters)

	address, err := madock8s.getServiceEndpoint(madock8sID)
	if err != nil {
		err = errors.Wrapf(err, "cannot handle change for %v", madock8sID)
		logrus.Error(err)
		return
	}
	logrus.Infof("service cluster ip: %v", address)
	madock8s.notifyChangeExporters(exporters, madock8sID, address)
	adapter.CreateDauxGenerateJob()
}

func (madock8s *Madock8s) notifyChangeExporters(exporters []notifier.Exporter, serviceName string, address string) {
	for _, exporter := range exporters {
		notifier.NotifyExporter(exporter, address, serviceName, madock8s.Namespace)
	}
}

func getMadock8sID(annotations map[string]string) string {
	return annotations["madock8s"]
}

func extractExporterNameAndValueName(annotation string) (string, string) {
	config := strings.Replace(annotation, "madock8s.exporter/", "", 1) // exporter.value
	slice := strings.Split(config, ".")
	name := slice[0]
	valueName := "value"
	if len(slice) > 1 {
		valueName = slice[1]
	}
	return name, valueName
}

func (madock8s *Madock8s) searchMadock8sExporter(annotations map[string]string) []notifier.Exporter {
	set := make(map[string]notifier.Exporter)
	for key, value := range annotations {
		if strings.HasPrefix(key, "madock8s.exporter/") {
			name, valueName := extractExporterNameAndValueName(key)

			exporter := notifier.Exporter{}
			// check if the exporter was already found and update it's config
			if e, ok := set[name]; ok {
				e.Configuration[valueName] = value
				exporter = e
			} else {
				config := map[string]string{valueName: value}
				exporter = notifier.Exporter{
					Name:          name,
					Configuration: config,
				}
			}
			set[name] = exporter
		}
	}

	exporters := []notifier.Exporter{}
	for _, exporter := range set {
		exporters = append(exporters, exporter)
	}
	return exporters
}

func (madock8s *Madock8s) resolveExporters(exporters []notifier.Exporter) []notifier.Exporter {
	resolvedExporters := make([]notifier.Exporter, len(exporters))
	for idx, exporter := range exporters {
		exporterClusterIP, err := madock8s.getServiceEndpoint(exporter.Name)
		if err != nil {
			err = errors.Wrapf(err, "cannot find exporter %v", exporter.Name)
			logrus.Error(err)
		} else {
			logrus.Infof("exporter ip for %v is %v", exporter.Name, exporterClusterIP)
			exporter.ClusterIP = exporterClusterIP
			resolvedExporters[idx] = exporter
		}
	}
	return resolvedExporters
}

func (madock8s *Madock8s) getServiceEndpoint(madock8sID string) (string, error) {
	services, err := madock8s.GetServices(madock8sID, madock8s.Namespace)
	if err != nil {
		err = errors.Wrapf(err, "Services not found: %v", err)
		return "", err
	}

	if len(services.Items) == 0 {
		return "", errors.New("No service endpoint found")
	}
	clusterIP := services.Items[0].Spec.ClusterIP + ":" + strconv.Itoa(int(services.Items[0].Spec.Ports[0].Port))
	return clusterIP, nil
}
