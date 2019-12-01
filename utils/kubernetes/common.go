package kubernetes

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getKubeClientSet returns a kubernetes client set which can be used to connect to kubernetes cluster
func getKubeClient(incluster bool, filepath string) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var clientset *kubernetes.Clientset
	var err error

	if incluster == true {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", filepath)
	}

	if err != nil {
		return nil, err
	}

	clientset, err = kubernetes.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getObject(client *kubernetes.Clientset, kubetest *kubeTest) (interface{}, error) {

	var listOptions *metav1.ListOptions

	if kubetest.Metadata.Label.Key != "" {
		listOptions = getListOptions(kubetest.Metadata.Label.Key, kubetest.Metadata.Label.Value)
	}

	switch kubetest.Metadata.Kind {

	case "Pod":
		if kubetest.Metadata.Name != "" {
			return getPod(client, kubetest.Metadata.Name, kubetest.Metadata.Namespace)
		}
		return listPods(client, kubetest.Metadata.Namespace, listOptions)

	case "Deployment":
		return getDeployment(client, kubetest.Metadata.Name, kubetest.Metadata.Namespace)

	default:
		return nil, fmt.Errorf("Cannot detect object type")

	}
}

func decodeTestFile(filePath string) (*kubeTest, error) {
	kubeTest := &kubeTest{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &kubeTest)
	if err != nil {
		return nil, err
	}

	return kubeTest, nil
}

func getListOptions(key, value string) *metav1.ListOptions {
	return &metav1.ListOptions{LabelSelector: key + "=" + value}
}

func getSlice(t interface{}) []interface{} {
	var slicedInterface []interface{}

	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			slicedInterface = append(slicedInterface, s.Index(i).Interface())
		}

	default:
		slicedInterface = append(slicedInterface, t)
	}

	return slicedInterface
}

// ApplyFile mimics kubectl apply -f. Takes in a path to a file and applies that object to the cluster and returns the applied object.
func applyAction(client *kubernetes.Clientset, pathToFile, action string) (*kubeMetadata, interface{}, error) {

	decode := scheme.Codecs.UniversalDeserializer().Decode

	bytes, err := ioutil.ReadFile(pathToFile)
	objectMetadata := &kubeMetadata{}
	var object interface{}

	if err != nil {
		log.Fatal(err, "Error reading the test file.")
	}

	obj, _, err := decode(bytes, nil, nil)

	if err != nil {
		log.Fatal(err, "Error decoding KubeTest object.")
	}

	err = yaml.Unmarshal(bytes, &objectMetadata)

	if err != nil {
		log.Fatal(err, "Error while unmarshalling KubeTest Metadata.")
	}

	log.InfoWithFields(map[string]interface{}{
		"action":    action,
		"name":      objectMetadata.Metadata.Name,
		"namespace": objectMetadata.Metadata.Namespace,
	}, "Applying action to file")

	switch obj.(type) {
	case *appsv1.Deployment:
		deploy := obj.(*appsv1.Deployment)
		if action == "CREATE" {
			object, err = client.AppsV1().Deployments(objectMetadata.Metadata.Namespace).Create(deploy)
		} else if action == "UPDATE" {
			object, err = client.AppsV1().Deployments(objectMetadata.Metadata.Namespace).Update(deploy)
		} else {
			err = client.AppsV1().Deployments(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1.Pod:
		pod := obj.(*v1.Pod)
		if action == "CREATE" {
			object, err = client.CoreV1().Pods(objectMetadata.Metadata.Namespace).Create(pod)
		} else if action == "UPDATE" {
			object, err = client.CoreV1().Pods(objectMetadata.Metadata.Namespace).Update(pod)
		} else {
			err = client.CoreV1().Pods(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1.Service:
		service := obj.(*v1.Service)
		if action == "CREATE" {
			object, err = client.CoreV1().Services(objectMetadata.Metadata.Namespace).Create(service)
		} else if action == "UPDATE" {
			object, err = client.CoreV1().Services(objectMetadata.Metadata.Namespace).Update(service)
		} else {
			err = client.CoreV1().Services(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1beta1.Ingress:
		ingress := obj.(*v1beta1.Ingress)
		if action == "CREATE" {
			object, err = client.ExtensionsV1beta1().Ingresses(objectMetadata.Metadata.Namespace).Create(ingress)
		} else if action == "UPDATE" {
			object, err = client.ExtensionsV1beta1().Ingresses(objectMetadata.Metadata.Namespace).Update(ingress)
		} else {
			err = client.ExtensionsV1beta1().Ingresses(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1beta1.DaemonSet:
		ds := obj.(*v1beta1.DaemonSet)
		if action == "CREATE" {
			object, err = client.ExtensionsV1beta1().DaemonSets(objectMetadata.Metadata.Namespace).Create(ds)
		} else if action == "UPDATE" {
			object, err = client.ExtensionsV1beta1().DaemonSets(objectMetadata.Metadata.Namespace).Update(ds)
		} else {
			err = client.ExtensionsV1beta1().DaemonSets(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	default:
		object, err = nil, fmt.Errorf("ApplyAction for kind is not implemented")
	}

	return objectMetadata, object, err
}
