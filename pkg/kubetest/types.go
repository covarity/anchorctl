package kubetest

import (
	"k8s.io/client-go/kubernetes"
)

type kubeTest struct {
	ApiVersion   string `yaml:"apiVersion"`
	Kind         string
	Metadata     metadata
	Spec         kubeTestSpec
	testFilePath string
}

type kubeTestSpec struct {
	Lifecycle lifecycle `yaml:"lifecycle"`
	Tests     []test
}

type test struct {
	Resource resource `yaml:"resource"`
	Type     string
	Spec     map[string]interface{}
}

type metadata struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

type kubeTester interface {
	test(res *resource) (bool, error)
}

type jsonTest struct {
	JsonPath string
	Value    string
	client   *kubernetes.Clientset
}

type validationTest struct {
	ContainsResponse string `yaml:"containsResponse"`
	client           *kubernetes.Clientset
}

type mutationTest struct {
	JsonPath string
	Value    string
	client   *kubernetes.Clientset
}
