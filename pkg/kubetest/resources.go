package kubetest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

type resource struct {
	ObjectRef objectRef `yaml:"objectRef"`
	Manifest  manifest  `yaml:"manifest"`
}

type objectRef struct {
	Type string        `yaml:"type"`
	Spec objectRefSpec `yaml:"spec"`
}

type lifecycle struct {
	PostStart []manifest `yaml:"postStart"`
	PreStop   []manifest `yaml:"preStop"`
}

func executeLifecycle(manifests []manifest, client *kubernetes.Clientset) {
	for _, i := range manifests {
		_, _, err := i.apply(client)
		if err != nil {
			log.Fatal(err, "Failed Post Start steps")
		}
	}
}

func (ob objectRef) valid() bool {
	if ob.Type == "" || ob.Spec.Kind == "" || ob.Spec.Namespace == "" ||
		ob.Spec.LabelValue == "" || ob.Spec.LabelKey == "" {

		log.WarnWithFields(map[string]interface{}{
			"resource": "objectRef",
			"expected": "Resource ObjectRef type, kind, namespace, label value and label key should be specified",
			"got":      "Type: " + ob.Type + " Kind: " + ob.Spec.Kind + " Namespace: " + ob.Spec.Namespace + " LabelKey: " + ob.Spec.LabelKey + " LabelValue: " + ob.Spec.LabelValue,
		}, "Failed getting the resource to apply.")

		return false
	}
	return true
}

func (ob objectRef) getObject(client *kubernetes.Clientset) (interface{}, error) {

	if valid := ob.valid(); valid != true {
		return nil, fmt.Errorf("AssertJSONPath object ref is invalid")
	}

	listOptions := getListOptions(ob.Spec.LabelKey, ob.Spec.LabelValue)

	if ob.Type == "Resource" {
		switch ob.Spec.Kind {

		case "Pod":
			return listPods(client, ob.Spec.Namespace, listOptions)

		case "ConfigMap":
			return listConfigMaps(client, ob.Spec.Namespace, listOptions)

		default:
			return nil, fmt.Errorf("Cannot detect object type")
		}
	} else {
		return nil, fmt.Errorf("Unknown objectRef type %s", ob.Type)
	}
}

type objectRefSpec struct {
	Kind       string `yaml:"kind"`
	Namespace  string `yaml:"namespace"`
	LabelKey   string `yaml:"labelKey"`
	LabelValue string `yaml:"labelValue"`
}

type manifest struct {
	Path   string `yaml:"path"`
	Action string `yaml:"action"`
}

func (mf manifest) valid() bool {
	if mf.Path == "" || mf.Action == "" {
		log.WarnWithFields(map[string]interface{}{
			"resource": "manifest",
			"expected": "Resource Manifest path and action should be specified",
			"got":      "Path: " + mf.Path + " Action: " + mf.Action,
		}, "Failed getting the resource to apply.")
		return false
	}
	return true
}

// ApplyFile mimics kubectl apply -f. Takes in a path to a file and applies that object to the cluster and returns the applied object.
func (mf manifest) apply(client *kubernetes.Clientset) (*kubeMetadata, interface{}, error) {

	if valid := mf.valid(); valid != true {
		return nil, nil, fmt.Errorf("Invalid Manifest to apply")
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	objectMetadata := &kubeMetadata{}
	var object interface{}

	bytes, err := ioutil.ReadFile(mf.Path)
	if err != nil {
		log.Error(err, "Error reading the "+mf.Path+" file.")
		return nil, nil, err
	}

	obj, _, err := decode(bytes, nil, nil)
	if err != nil {
		log.Error(err, "Error decoding bytes to kube object.")
		return nil, nil, err
	}

	err = yaml.Unmarshal(bytes, &objectMetadata)
	if err != nil {
		log.Error(err, "Error while unmarshalling KubeTest Metadata.")
	}

	log.InfoWithFields(map[string]interface{}{
		"action":    mf.Action,
		"name":      objectMetadata.Metadata.Name,
		"namespace": objectMetadata.Metadata.Namespace,
	}, "Applying action to file")

	switch obj.(type) {
	case *appsv1.Deployment:
		deploy := obj.(*appsv1.Deployment)
		if mf.Action == "CREATE" {
			object, err = client.AppsV1().Deployments(objectMetadata.Metadata.Namespace).Create(deploy)
		} else if mf.Action == "UPDATE" {
			object, err = client.AppsV1().Deployments(objectMetadata.Metadata.Namespace).Update(deploy)
		} else {
			err = client.AppsV1().Deployments(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1.Pod:
		pod := obj.(*v1.Pod)
		if mf.Action == "CREATE" {
			object, err = client.CoreV1().Pods(objectMetadata.Metadata.Namespace).Create(pod)
		} else if mf.Action == "UPDATE" {
			object, err = client.CoreV1().Pods(objectMetadata.Metadata.Namespace).Update(pod)
		} else {
			err = client.CoreV1().Pods(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1.Service:
		service := obj.(*v1.Service)
		if mf.Action == "CREATE" {
			object, err = client.CoreV1().Services(objectMetadata.Metadata.Namespace).Create(service)
		} else if mf.Action == "UPDATE" {
			object, err = client.CoreV1().Services(objectMetadata.Metadata.Namespace).Update(service)
		} else {
			err = client.CoreV1().Services(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1beta1.Ingress:
		ingress := obj.(*v1beta1.Ingress)
		if mf.Action == "CREATE" {
			object, err = client.ExtensionsV1beta1().Ingresses(objectMetadata.Metadata.Namespace).Create(ingress)
		} else if mf.Action == "UPDATE" {
			object, err = client.ExtensionsV1beta1().Ingresses(objectMetadata.Metadata.Namespace).Update(ingress)
		} else {
			err = client.ExtensionsV1beta1().Ingresses(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1beta1.DaemonSet:
		ds := obj.(*v1beta1.DaemonSet)
		if mf.Action == "CREATE" {
			object, err = client.ExtensionsV1beta1().DaemonSets(objectMetadata.Metadata.Namespace).Create(ds)
		} else if mf.Action == "UPDATE" {
			object, err = client.ExtensionsV1beta1().DaemonSets(objectMetadata.Metadata.Namespace).Update(ds)
		} else {
			err = client.ExtensionsV1beta1().DaemonSets(objectMetadata.Metadata.Namespace).Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	case *v1.Namespace:
		ns := obj.(*v1.Namespace)
		if mf.Action == "CREATE" {
			object, err = client.CoreV1().Namespaces().Create(ns)
		} else if mf.Action == "UPDATE" {
			object, err = client.CoreV1().Namespaces().Update(ns)
		} else {
			err = client.CoreV1().Namespaces().Delete(objectMetadata.Metadata.Name, &metav1.DeleteOptions{})
		}
	default:
		object, err = nil, fmt.Errorf("ApplyAction for kind is not implemented")
	}

	if err != nil {
		log.WarnWithFields(map[string]interface{}{
			"Stage":  "PostStart",
			"Path":   mf.Path,
			"Action": mf.Action,
			"Error":  err.Error(),
		}, "Apply Manifest error.")
	}

	return objectMetadata, object, err
}
