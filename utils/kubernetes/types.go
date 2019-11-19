package kubernetes

type kubeTest struct{
	Kind string
	Metadata metadata
	Tests []map[string]string
}

type metadata struct {
	Name string
	Namespace string
	Label label
}

type label struct {
	Key string
	Value string
}
