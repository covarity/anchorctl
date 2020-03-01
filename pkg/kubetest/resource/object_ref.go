package resource

import (
	"anchorctl/pkg/logging"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"os/exec"
	"path/filepath"
)

var log *logging.Logger
var TestFilePath string

type Resource struct {
	ObjectRef objectRef `yaml:"objectRef"`
	Manifest  Manifest  `yaml:"Manifest"`
}

type objectRef struct {
	Type string        `yaml:"type"`
	Spec objectRefSpec `yaml:"spec"`
}

type objectRefSpec struct {
	Kind      string `yaml:"kind"`
	Namespace string `yaml:"namespace"`
	Container string `yaml:"container,omitempty"`
	Labels    map[string]string
}

type Manifest struct {
	Path   string `yaml:"path"`
	Action string `yaml:"action"`
}

func (ob objectRef) GetObject(client *kubernetes.Clientset) ([]runtime.Object, error) {

	if valid := ob.valid(); !valid {
		return nil, fmt.Errorf("assertJSONPath object ref is invalid")
	}

	if ob.Type != "Resource" {
		return nil, fmt.Errorf("unknown objectRef type %s", ob.Type)
	}

	api, isNamespaced, e := ob.getAPI(client)
	if e != nil {
		return nil, e
	}

	req := api.Get().Resource(pluralise(ob.Spec.Kind))

	if isNamespaced {
		req = req.Namespace(ob.Spec.Namespace)
	}

	for k, v := range ob.Spec.Labels {
		req.Param("labelSelector", k+"="+v)
	}

	res, err := req.Do().Get()

	if err != nil {
		return nil, err
	}

	return meta.ExtractList(res)
}

func (ob objectRef) getAPI(client *kubernetes.Clientset) (rest.Interface, bool, error) {
	var api rest.Interface
	var isNamespaced = true

	switch ob.Spec.Kind {
	case "DaemonSet", "Deployment", "ReplicaSet", "StatefulSet":
		api = client.AppsV1().RESTClient()

	case "HorizontalPodAutoscaler":
		api = client.AutoscalingV1().RESTClient()

	case "Job":
		api = client.BatchV1().RESTClient()

	case "ComponentStatus", "Namespace", "Node", "PersistentVolume":
		isNamespaced = false
		api = client.CoreV1().RESTClient()

	case "ConfigMap", "Endpoint", "Event", "LimitRange", "PersistentVolumeClaim", "Pod",
		"ResourceQuota", "Secrets", "Services", "ServiceAccount":
		api = client.CoreV1().RESTClient()

	case "NetworkPolicy":
		api = client.NetworkingV1().RESTClient()

	case "PodDisruptionBudget", "PodSecurityPolicies":
		api = client.PolicyV1beta1().RESTClient()

	case "ClusterRole", "ClusterRoleBinding":
		isNamespaced = false
		api = client.RbacV1().RESTClient()

	case "Role", "RoleBinding":
		api = client.RbacV1().RESTClient()

	case "PriorityClass":
		api = client.SchedulingV1().RESTClient()

	default:
		return nil, false, fmt.Errorf("api not found for the object")
	}
	return api, isNamespaced, nil
}

func pluralise(str string) string {

	exceptions := make(map[string]string)

	pluralise := namer.NewAllLowercasePluralNamer(exceptions)

	pluralType := types.Type{
		Name: types.Name{Name: str},
	}

	return pluralise.Name(&pluralType)
}

// ApplyFile mimics kubectl apply -f. Takes in a path to a file and applies that object to the
// cluster and returns the applied object.
func (mf Manifest) Apply(expectError bool) (*objectRef, error) {

	if valid := mf.valid(); !valid {
		return nil, fmt.Errorf("invalid Manifest to apply")
	}

	var filePath string
	if TestFilePath != "" {
		filePath = filepath.Clean(filepath.Join(filepath.Dir(TestFilePath), mf.Path))
	} else {
		filePath = filepath.Clean(mf.Path)
	}

	objRef, err := mf.getObjectref(filePath)
	if err != nil {
		log.Error(err, "Error Applying action to file")
		return nil, err
	}

	log.InfoWithFields(map[string]interface{}{
		"Action": mf.Action,
		"Path":   filePath,
	}, "Applying action to file")

	var out []byte

	if mf.Action == "CREATE" {
		// #nosec
		out, err = exec.Command("kubectl", "apply", "-f", filePath).CombinedOutput()
	} else {
		// #nosec
		out, err = exec.Command("kubectl", "delete", "-f", filePath).CombinedOutput()
	}

	if err != nil {
		if !expectError {
			log.WarnWithFields(map[string]interface{}{
				"Path":   filePath,
				"Action": mf.Action,
				"Error":  err.Error(),
			}, "Apply Manifest error.")
		}

		applyError := fmt.Errorf(string(out))

		return nil, applyError
	}

	log.InfoWithFields(map[string]interface{}{
		"action":   mf.Action,
		"filepath": filePath,
	}, string(out))

	return objRef, err
}

func (mf Manifest) getObjectref(filePath string) (*objectRef, error) {
	ymlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err, "Error reading the "+mf.Path+" file.")
		return nil, err
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(ymlFile))
	if err != nil {
		log.Error(err, "Error reading Test file.")
	}

	objRef := &objectRef{
		Type: "Resource",
		Spec: objectRefSpec{
			Kind:      viper.GetString("kind"),
			Namespace: viper.GetString("metadata.namespace"),
			Labels:    viper.GetStringMapString("metadata.labels"),
		},
	}
	return objRef, err
}

func listPods(client *kubernetes.Clientset, namespace string, listOptions metav1.ListOptions) ([]v1.Pod, error) {
	pods, err := client.CoreV1().Pods(namespace).List(listOptions)
	return pods.Items, err
}

func (ob objectRef) valid() bool {
	if ob.Type == "" || ob.Spec.Kind == "" || ob.Spec.Namespace == "" ||
		ob.Spec.Labels == nil {

		log.WarnWithFields(map[string]interface{}{
			"Resource": "objectRef",
			"expected": "Resource ObjectRef type, kind, namespace, label value and label key should be specified",
			"got":      "Type: " + ob.Type + " Kind: " + ob.Spec.Kind + " Namespace: " + ob.Spec.Namespace,
		}, "Failed getting the Resource to apply.")

		return false
	}

	return true
}

func (mf Manifest) valid() bool {
	if mf.Path == "" || mf.Action == "" {
		log.WarnWithFields(map[string]interface{}{
			"Resource": "Manifest",
			"expected": "Resource Manifest path and action should be specified",
			"got":      "Path: " + mf.Path + " Action: " + mf.Action,
		}, "Failed getting the Resource to apply.")

		return false
	}

	return true
}

func ExecuteLifecycle(manifests []Manifest) {
	for _, i := range manifests {
		_, err := i.Apply(false)
		if err != nil {
			log.Fatal(err, "Failed Lifecycle steps")
		}
	}
}
