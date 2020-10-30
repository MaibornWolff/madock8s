package adapter

import (
	"github.com/MaibornWolff/maDocK8s/core/controller/utils"
	"github.com/sirupsen/logrus"
	batch "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespace = utils.GetEnvTargetNamespace()
var jobName = "daux-generator"

var yaml = `
apiVersion: batch/v1
kind: Job
metadata:
  name: daux-generator
spec:
  backoffLimit: 4
  template:
    spec:
      containers:
        - name: daux
          image: daux/daux.io:latest
          imagePullPolicy: IfNotPresent
          command:
            - daux
            - generate
          volumeMounts:
            - mountPath: /build/docs
              name: docs-volume
            - mountPath: /build/static
              name: static-volume

        - name: index
          imagePullPolicy: IfNotPresent
          image: busybox
          securityContext:
           runAsUser: 0
          volumeMounts:
          - mountPath: /var/www
            name: static-volume
          - mountPath: "/var/index.html"
            name: index-page
            subPath: "index.html"
          command:
          - sh
          - "-c"
          args:
          - "mv /var/index.html /var/www/index.html; ls /var/www"

      restartPolicy: Never

      volumes:
        - name: docs-volume
          persistentVolumeClaim:
            claimName: daux-generator-in
        - name: static-volume
          persistentVolumeClaim:
            claimName: daux-generator-out
        - name: index-page
          configMap:
            name: index-page

`

func CreateDauxGenerateJob() {
	clientset := utils.GetClient()

	var period int64 = 0
	delOpts := metav1.DeleteOptions{
		GracePeriodSeconds: &period,
	}
	err := clientset.BatchV1().Jobs(namespace).Delete(jobName, &delOpts)
	if err != nil {
		logrus.Errorf("could not delete job: %v", err)
	}
	logrus.Infof("Job %s is deleted.", jobName)

	time.Sleep(5 * time.Second)
	createJob(clientset)
}

func parseYaml() (*batch.Job, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(yaml), nil, nil)
	if err != nil {
		return nil, err
	}
	job := obj.(*batch.Job)
	logrus.Infof("%#v\n", utils.SimplifyLog(job))

	return job, nil
}

func createJob(clientset kubernetes.Interface) {
	job, err := parseYaml()
	if err != nil {
		logrus.Fatalf("could not parse Job %v", err)
	}

	createdJob, err := clientset.BatchV1().Jobs(namespace).Create(job)
	if err != nil {
		logrus.Fatalf("could not create job: %v", err)
	}
	logrus.Infof("job created: %s", createdJob.Name)
}
