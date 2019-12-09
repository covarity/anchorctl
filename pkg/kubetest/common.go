package kubetest

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func getObject(client *kubernetes.Clientset, ref *objectRef) (interface{}, error) {

	if ref.Type == "File" {
		return nil, nil
	}

	var listOptions *metav1.ListOptions

	if ref.Spec.LabelKey != "" {
		listOptions = getListOptions(ref.Spec.LabelKey, ref.Spec.LabelValue)
	}

	switch ref.Type {

	case "Pod":
		return listPods(client, ref.Spec.Namespace, listOptions)

	case "ConfigMap":
		return listConfigMaps(client, ref.Spec.Namespace, listOptions)

	default:
		return nil, fmt.Errorf("Cannot detect object type")

	}
}

func decodeTestFile(client *kubernetes.Clientset, filePath string) (*kubeTest, error) {
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

func getListOptions(key, value string) *metav1.ListOptions {
	return &metav1.ListOptions{LabelSelector: key + "=" + value}
}

func getSlice(t interface{}) []interface{} {
	var slicedInterface []interface{}

	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			slicedInterface = append(slicedInterface, s.Index(i).Interface())
		}

	default:
		slicedInterface = append(slicedInterface, t)
	}

	return slicedInterface
}
