package kubetest

import (
	"k8s.io/client-go/kubernetes"
)

type kubeTest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   metadata
	Spec       kubeTestSpec
}

type kubeTestSpec struct {
	Tests []test
}

type test struct {
	ObjectRef objectRef `yaml:"objectRef"`
	Type      string
	Spec      map[string]interface{}
	Assertion kubeTester
}

type metadata struct {
	Name      string
	Namespace string
	Label     map[string]string
}

type kubeMetadata struct {
	Metadata metadata
}

type objectRef struct {
	Type string        `yaml:"type"`
	Spec objectRefSpec `yaml:"spec"`
}

type objectRefSpec struct {
	Path       string `yaml:"path"`
	Action     string `yaml:"action"`
	Kind       string `yaml:"kind"`
	Namespace  string `yaml:"namespace"`
	LabelKey   string `yaml:"labelKey"`
	LabelValue string `yaml:"labelValue"`
}

type kubeTester interface {
	test(objectRef objectRef) (bool, error)
}

type jsonTest struct {
	JsonPath string
	Value    string
	client   *kubernetes.Clientset
}

type validationTest struct {
	ExpectedError string
	client        *kubernetes.Clientset
}

type mutationTest struct {
	JsonPath string
	Value    string
	client   *kubernetes.Clientset
}
