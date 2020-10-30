package app

import (
	"testing"

	"github.com/MaibornWolff/maDocK8s/core/controller/adapter"
	"github.com/MaibornWolff/maDocK8s/core/controller/notifier"
	. "github.com/onsi/gomega"
	api_v1 "k8s.io/api/core/v1"
)

func Test_Should_Return_2_Exporter_From_Annotations(t *testing.T) {
	Ω := NewGomegaWithT(t)
	madock8s := &Madock8s{}
	annotations := map[string]string{
		"madock8s.exporter/gitlabExporter.baseurl": "http://mygitlab.de",
		"madock8s.exporter/prometheusExporter":     "true",
		"nomadock8sexporter/noExporter":            "true",
	}

	exporters := madock8s.searchMadock8sExporter(annotations)

	expExporter1 := notifier.Exporter{
		Name:          "gitlabExporter",
		Configuration: map[string]string{"baseurl": "http://mygitlab.de"},
	}

	expExporter2 := notifier.Exporter{
		Name:          "prometheusExporter",
		Configuration: map[string]string{"value": "true"},
	}

	Ω.Expect(exporters).Should(ConsistOf(expExporter1, expExporter2))
}

func Test_Should_Return_An_Empty_Exporter_List_With_No_Exporter_Annotations(t *testing.T) {
	Ω := NewGomegaWithT(t)
	madock8s := &Madock8s{}

	annotations := map[string]string{
		"madock8s":                      "true",
		"nomadock8sexporter/noExporter": "true",
	}

	exporters := madock8s.searchMadock8sExporter(annotations)

	Ω.Expect(exporters).Should(BeEmpty())
}

func Test_Should_Get_ClusterIP_From_Service_List(t *testing.T) {
	Ω := NewGomegaWithT(t)

	serviceList := &api_v1.ServiceList{
		Items: []api_v1.Service{
			api_v1.Service{
				Spec: api_v1.ServiceSpec{
					ClusterIP: "0.0.0.0",
					Ports: []api_v1.ServicePort{
						{Name: "http", Port: 80},
					},
				},
			},
		},
	}
	mock := adapter.ServicesMock{
		ServiceList: serviceList,
		ID:          "myService",
	}

	madock8s := &Madock8s{
		GetServices: mock.GetServices,
	}

	clusterIP, _ := madock8s.getServiceEndpoint("myService")

	Ω.Expect(clusterIP).To(Equal("0.0.0.0:80"))
}

func Test_Should_Return_An_Error_With_Empty_List_Message(t *testing.T) {
	Ω := NewGomegaWithT(t)

	serviceList := &api_v1.ServiceList{
		Items: []api_v1.Service{},
	}

	mock := adapter.ServicesMock{
		ServiceList: serviceList,
		ID:          "myService",
	}

	madock8s := &Madock8s{
		GetServices: mock.GetServices,
	}

	clusterIP, err := madock8s.getServiceEndpoint("myService")

	Ω.Expect(clusterIP).Should(BeEmpty())
	Ω.Expect(err).Should(MatchError("No service endpoint found"))
}

func Test_Should_Take_First_ClusterIP_From_Service_List(t *testing.T) {
	Ω := NewGomegaWithT(t)

	service1 := api_v1.Service{
		Spec: api_v1.ServiceSpec{
			ClusterIP: "1.0.0.0",
			Ports: []api_v1.ServicePort{
				{Name: "http", Port: 80},
			},
		},
	}

	service2 := api_v1.Service{
		Spec: api_v1.ServiceSpec{
			ClusterIP: "2.0.0.0",
			Ports: []api_v1.ServicePort{
				{Name: "http", Port: 80},
			},
		},
	}

	serviceList := &api_v1.ServiceList{
		Items: []api_v1.Service{
			service1,
			service2,
		},
	}
	mock := adapter.ServicesMock{
		ServiceList: serviceList,
		ID:          "myService",
	}
	madock8s := &Madock8s{
		GetServices: mock.GetServices,
	}

	clusterIP, _ := madock8s.getServiceEndpoint("myService")

	Ω.Expect(clusterIP).To(Equal("1.0.0.0:80"))
}
