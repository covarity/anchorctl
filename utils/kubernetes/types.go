package kubernetes

type kubeTest struct {
	Metadata metadata
	Tests    []map[string]string
}

type metadata struct {
	Kind      string
	Name      string
	Namespace string
	Label     label
}

type label struct {
	Key   string
	Value string
}

type kubeMetadata struct {
	Metadata metadata
}
