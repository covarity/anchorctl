package kubetest

import (
	"anchorctl/pkg/logging"
	"anchorctl/pkg/resultaggregator"
	"k8s.io/client-go/kubernetes"
)

type KubeTest struct {
	Lifecycle lifecycle `yaml:"lifecycle"`
	Tests     []test
	Opts	  Options
}

type test struct {
	Resource resource `yaml:"resource"`
	Type     string
	Spec     map[string]interface{}
}

type kubeTester interface {
	test(res *resource) *resultaggregator.TestRun
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

type Options struct {
	Incluster bool
	Kubeconfig string
	TestFilepath string
	Logger  *logging.Logger
}

type resource struct {
	ObjectRef objectRef `yaml:"objectRef"`
	Manifest  manifest  `yaml:"manifest"`
}

type objectRef struct {
	Type string        `yaml:"type"`
	Spec objectRefSpec `yaml:"spec"`
}

type lifecycle struct {
	PostStart []manifest `yaml:"postStart"`
	PreStop   []manifest `yaml:"preStop"`
}

type objectRefSpec struct {
	Kind      string `yaml:"kind"`
	Namespace string `yaml:"namespace"`
	Labels    map[string]string
}

type manifest struct {
	Path   string `yaml:"path"`
	Action string `yaml:"action"`
}