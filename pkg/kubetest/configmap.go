package kubetest

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPod
func getConfigMap(client *kubernetes.Clientset, name, namespace string) ([]v1.ConfigMap, error) {
	configmaps, err := client.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	return []v1.ConfigMap{*configmaps}, err
}

func listConfigMaps(client *kubernetes.Clientset, namespace string, listOptions *metav1.ListOptions) ([]v1.ConfigMap, error) {
	configmaps, err := client.CoreV1().ConfigMaps(namespace).List(*listOptions)
	return configmaps.Items, err
}
