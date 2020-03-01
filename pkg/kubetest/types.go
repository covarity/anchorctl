package kubetest

import (
	"anchorctl/pkg/kubetest/resource"
	"anchorctl/pkg/logging"
)

type KubeTest struct {
	Lifecycle lifecycle `yaml:"lifecycle"`
	Tests     []test
	Opts      Options
}

type test struct {
	Resource resource.Resource `yaml:"Resource"`
	Type     TestType
	Spec     map[string]interface{}
}

type Options struct {
	Incluster    bool
	Kubeconfig   string
	TestFilepath string
	Logger       *logging.Logger
}

type lifecycle struct {
	PostStart []resource.Manifest `yaml:"postStart"`
	PreStop   []resource.Manifest `yaml:"preStop"`
}
