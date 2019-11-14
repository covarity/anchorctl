package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPod
func GetPod(client *kubernetes.Clientset, name, namespace string) (*v1.Pod, error) {
	return client.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
}

func ListPods(client *kubernetes.Clientset, namespace string, listOptions metav1.ListOptions) (*v1.PodList, error) {
	return client.CoreV1().Pods(namespace).List(listOptions)
}

func ListPodsByLabel(client *kubernetes.Clientset, namespace, key, value string) (*v1.PodList, error) {
	listOptions := metav1.ListOptions{LabelSelector: key + "=" + value}
	return ListPods(client, namespace, listOptions)
}

func ListPodNamesByLabel(client *kubernetes.Clientset, namespace, key, value string) ([]string, error) {
	pods, err := ListPodsByLabel(client, namespace, key, value)

	if err != nil {
		return nil, err
	}

	var names []string

	for _, i := range pods.Items {
		names = append(names, i.Name)
	}

	return names, err
}

