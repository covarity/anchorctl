package kubetest

import (
	"k8s.io/client-go/kubernetes"
)

type kubeTest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   metadata
	Spec       kubeTestSpec
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
	JSONPath string
	Value    string
	client   *kubernetes.Clientset
}

type validationTest struct {
	ContainsResponse string `yaml:"containsResponse"`
	client           *kubernetes.Clientset
}

type mutationTest struct {
	JSONPath string
	Value    string
	client   *kubernetes.Clientset
}
