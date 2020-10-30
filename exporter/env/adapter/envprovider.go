package adapter

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ContainerVarsMap map[string][]EnvVar

type EnvVar struct {
	Key           string
	ExtSourceType string
	ExtSourceName string
	ExtSourceKey  string
	Value         string
}

func GetEnvironmentVars(deploymentName string, namespace string) (ContainerVarsMap, error) {

	containerVarsMap := make(ContainerVarsMap)

	kubeclient := GetClient()

	deployment, err := kubeclient.AppsV1beta1().Deployments(namespace).Get(deploymentName, meta_v1.GetOptions{})
	if err != nil {
		return containerVarsMap, errors.Wrap(err, "failed GET for deployment")
	}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		containerVarsMap[container.Name] = getVarsFromContainer(kubeclient, container, namespace)
	}
	return containerVarsMap, nil
}

func getVarsFromContainer(kubeClient kubernetes.Interface, container v1.Container, namespace string) []EnvVar {
	envVarsArray := []EnvVar{}

	for _, v1EnvVar := range container.Env {
		resultEnvVar := EnvVar{
			Key:   v1EnvVar.Name,
			Value: v1EnvVar.Value,
		}

		if v1EnvVar.ValueFrom != nil {
			var err error
			resultEnvVar, err = getReferencedEnvVar(resultEnvVar, kubeClient, *v1EnvVar.ValueFrom, namespace)
			if err != nil {
				logrus.Error(errors.Wrap(err, "failed extracting referenced value"))
			}
		}
		envVarsArray = append(envVarsArray, resultEnvVar)
	}
	return envVarsArray
}

func getReferencedEnvVar(envVar EnvVar, kubeclient kubernetes.Interface, valueFrom v1.EnvVarSource, namespace string) (EnvVar, error) {

	if valueFrom.ConfigMapKeyRef != nil {
		return getVarFromConfigMap(envVar, kubeclient, *valueFrom.ConfigMapKeyRef, namespace)
	}

	if valueFrom.SecretKeyRef != nil {
		return getVarFromSecret(envVar, kubeclient, *valueFrom.SecretKeyRef, namespace)
	}

	if valueFrom.FieldRef != nil {
		return getVarFromFieldRef(envVar, kubeclient, *valueFrom.FieldRef)
	}

	if valueFrom.ResourceFieldRef != nil {
		return getVarFromResourceFieldRef(envVar, kubeclient, *valueFrom.ResourceFieldRef)
	}

	return EnvVar{}, errors.New("not supported value-referencing")
}

func getVarFromConfigMap(envVar EnvVar, kubeclient kubernetes.Interface, keyRef v1.ConfigMapKeySelector, namespace string) (EnvVar, error) {
	envVar.ExtSourceType = "configmap"
	envVar.ExtSourceName = keyRef.Name
	envVar.ExtSourceKey = keyRef.Key

	configMap, err := kubeclient.CoreV1().ConfigMaps(namespace).Get(envVar.ExtSourceName, meta_v1.GetOptions{})
	if err != nil {
		return envVar, errors.Wrap(err, "failed GET for configMap")
	}

	envVar.Value, err = getValueFrom(configMap.Data, envVar.ExtSourceKey)
	if err != nil {
		envVar.Value = "Value was not found in deployed configmap"
		return envVar, errors.Wrap(err, "")
	}

	return envVar, nil
}

func getVarFromFieldRef(envVar EnvVar, kubeclient kubernetes.Interface, keyRef v1.ObjectFieldSelector) (EnvVar, error) {
	envVar.ExtSourceType = "fieldRef"
	envVar.ExtSourceKey = keyRef.FieldPath
	envVar.Value = "Unable to provide. Check deployment"
	return envVar, nil
}

func getVarFromResourceFieldRef(envVar EnvVar, kubeClient kubernetes.Interface, keyRef v1.ResourceFieldSelector) (EnvVar, error) {
	envVar.ExtSourceType = "resourceFieldRef"
	envVar.ExtSourceName = keyRef.ContainerName
	envVar.ExtSourceKey = keyRef.Resource
	envVar.Value = "Unable to provide. Check deployment"
	return envVar, nil
}

func getVarFromSecret(envVar EnvVar, kubeclient kubernetes.Interface, keyRef v1.SecretKeySelector, namespace string) (EnvVar, error) {
	envVar.ExtSourceType = "secret"
	envVar.ExtSourceName = keyRef.Name
	envVar.ExtSourceKey = keyRef.Key

	secret, err := kubeclient.CoreV1().Secrets(namespace).Get(envVar.ExtSourceName, meta_v1.GetOptions{})
	if err != nil {
		return envVar, errors.Wrap(err, "failed GET for secret")
	}

	length := len(secret.Data[envVar.ExtSourceKey])
	// instead of a secret value, show asterisks
	envVar.Value = strings.Repeat("*", length)

	if length == 0 {
		envVar.Value = "Value was not found in deployed Secret"
		return envVar, errors.Errorf("value not found in %s for key %s", envVar.ExtSourceName, envVar.ExtSourceKey)
	}

	return envVar, nil
}

func getValueFrom(data map[string]string, key string) (string, error) {
	value := data[key]
	if value == "" {
		return "", errors.Errorf("value not found for key %v", key)
	}
	return value, nil
}
