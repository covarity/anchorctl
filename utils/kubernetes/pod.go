package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPod
func getPod(client *kubernetes.Clientset, name, namespace string) ([]v1.Pod, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	return []v1.Pod{*pod}, err
}

func listPods(client *kubernetes.Clientset, namespace string, listOptions *metav1.ListOptions) ([]v1.Pod, error) {
	pods, err := client.CoreV1().Pods(namespace).List(*listOptions)
	return pods.Items, err
}
