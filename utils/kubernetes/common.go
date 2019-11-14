package kubernetes

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClientSet returns a kubernetes client set which can be used to connect to kubernetes cluster
func GetKubeClient(incluster bool, filepath string) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var clientset *kubernetes.Clientset
	var err error

	if incluster == true {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", filepath)
	}

	if err != nil {
		return nil, err
	}

	clientset, err = kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetObject(client *kubernetes.Clientset, metadata *ObjectMetadata) (interface{}, error) {

	switch metadata.Kind {

	case "pod":
		return GetPod(client, metadata.Name, metadata.Namespace)

	case "deployment":
		return GetDeployment(client, metadata.Name, metadata.Namespace)

	default:
		return nil, fmt.Errorf("Cannot detect object type")

	}
}

func DecodeTestFile(filePath string) (*KubeTest, error) {
	kubeTest := &KubeTest{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &kubeTest)
	if err != nil {
		return nil, err
	}

	return kubeTest, nil
}
