package kubernetes

type KubeTest struct{
	Object ObjectMetadata
	Tests []map[string]string
}

type ObjectMetadata struct {
	Kind string
	Name string
	Namespace string
}
