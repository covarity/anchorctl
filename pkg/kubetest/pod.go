package kubetest

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
	log.InfoWithFields(map[string]interface{}{
		"kind":      "Pod",
		"Namespace": namespace,
		"labelkey":  listOptions.LabelSelector,
	}, "Retriving object")

	pods, err := client.CoreV1().Pods(namespace).List(*listOptions)
	if err != nil {
		log.Error(err, "Unable to retrieve object")
		return nil, err
	}

	return pods.Items, nil
}
