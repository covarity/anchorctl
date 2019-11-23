package kubernetes

type kubeTest struct{
	Object objectMetadata
	Tests []map[string]string
}

type objectMetadata struct {
	Kind string
	Name string
	Namespace string
	//file string
	Label label
}

type label struct {
	Key string
	Value string
}
