package common

import (
	"anchorctl/pkg/kubetest/resource"
	"anchorctl/pkg/resultaggregator"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

var RestKubeConfig *rest.Config

type KubeTester interface {
	Test(res *resource.Resource) *resultaggregator.TestRun
	Valid(res *resource.Resource) (bool, error)
}

func ValidTest(testType string, values map[string]string) (bool, error) {
	for k, v := range values {
		if len(v) <= 0 {
			err := fmt.Errorf("Test " + testType + " missing the following key:" + k)
			return false, err
		}
	}
	return true, nil
}

func MapToString(obj map[string]string, sep string) string {
	arr := []string{}

	for k, v := range obj {
		arr = append(arr, k+"="+v)
	}

	return strings.Join(arr, sep)
}

// getKubeClientSet returns a kubernetes client set which can be used to connect to kubernetes cluster
func GetKubeClient(incluster bool, filepath string) (*kubernetes.Clientset, error) {
	var config *rest.Config

	var clientset *kubernetes.Clientset

	var err error

	if incluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", filepath)
	}

	if err != nil {
		return nil, err
	}

	RestKubeConfig = config

	clientset, err = kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return clientset, nil
}
