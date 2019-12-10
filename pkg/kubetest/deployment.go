package kubetest

import (
	app "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getDeployment(client *kubernetes.Clientset, name, namespace string) (*app.Deployment, error) {
	return client.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}
