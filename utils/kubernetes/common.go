package kubernetes

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getKubeClientSet returns a kubernetes client set which can be used to connect to kubernetes cluster
func getKubeClient(incluster bool, filepath string) (*kubernetes.Clientset, error) {

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

func getObject(client *kubernetes.Clientset, metadata *objectMetadata) (interface{}, error) {

	switch metadata.Kind {

	case "pod":

		// Use Label if exists
		if metadata.Label.Key != "" {
			pods, err := listPodsByLabel(client, metadata.Namespace, metadata.Label.Key, metadata.Label.Value)
			if err != nil {
				return nil, fmt.Errorf("Cannot get pod list with key "+ metadata.Label.Key + " and value " + metadata.Label.Value, err)
			}

			if len(pods.Items) < 1 {
				return nil, fmt.Errorf("No pods with key "+ metadata.Label.Key + " and value " + metadata.Label.Value, err)
			}

			return pods.Items[0], nil
		}

		return getPod(client, metadata.Name, metadata.Namespace)

	case "deployment":
		return getDeployment(client, metadata.Name, metadata.Namespace)

	default:
		return nil, fmt.Errorf("Cannot detect object type")

	}
}

func decodeTestFile(filePath string) (*kubeTest, error) {
	kubeTest := &kubeTest{}

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
