package test

import (
	"anchorctl/pkg/resultaggregator"
)

type AnchorTest interface {
	Assert() *resultaggregator.TestResult
	//Verify(tester anchorTest) bool
}

type AnchorCRD struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       TestKind
	Metadata   metadata
	Spec       map[string]interface{}
}

type metadata struct {
	Name      string
	Namespace string
	Labels    map[string]string
}
