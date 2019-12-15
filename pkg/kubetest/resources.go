package kubetest

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"

	meta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"os/exec"
	"path/filepath"
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
		_, err := i.apply(client, false)
		if err != nil {
			log.Fatal(err, "Failed Lifecycle steps")
		}
	}
}

func (ob objectRef) valid() bool {
	if ob.Type == "" || ob.Spec.Kind == "" || ob.Spec.Namespace == "" ||
		ob.Spec.Labels == nil {

		log.WarnWithFields(map[string]interface{}{
			"resource": "objectRef",
			"expected": "Resource ObjectRef type, kind, namespace, label value and label key should be specified",
			"got":      "Type: " + ob.Type + " Kind: " + ob.Spec.Kind + " Namespace: " + ob.Spec.Namespace,
		}, "Failed getting the resource to apply.")

		return false
	}
	return true
}

func (ob objectRef) getObject(client *kubernetes.Clientset) ([]runtime.Object, error) {

	if valid := ob.valid(); !valid {
		return nil, fmt.Errorf("AssertJSONPath object ref is invalid")
	}

	if ob.Type != "Resource" {
		return nil, fmt.Errorf("Unknown objectRef type %s", ob.Type)
	}

	var api rest.Interface

	switch ob.Spec.Kind {
	case "DaemonSet", "Deployment", "ReplicaSet", "StatefulSet":
		api = client.AppsV1().RESTClient()

	case "HorizontalPodAutoscaler":
		api = client.AutoscalingV1().RESTClient()

	case "Job":
		api = client.BatchV1().RESTClient()

	case "ComponentStatus", "ConfigMap", "Endpoint", "Event", "LimitRange", "Namespace", "Node", "PersistentVolume", "PersistentVolumeClaim", "Pod", "ResourceQuota", "Secrets", "Services", "ServiceAccount":
		api = client.CoreV1().RESTClient()

	case "NetworkPolicy":
		api = client.NetworkingV1().RESTClient()

	case "PodDisruptionBudget", "PodSecurityPolicies":
		api = client.PolicyV1beta1().RESTClient()

	case "ClusterRole", "ClusterRoleBinding", "Role", "RoleBinding":
		api = client.RbacV1().RESTClient()

	case "PriorityClass":
		api = client.SchedulingV1().RESTClient()

	default:
		return nil, fmt.Errorf("Api not found for the object")
	}

	req := api.Get().Resource(pluralise(ob.Spec.Kind))

	for k, v := range ob.Spec.Labels {
		req.Param("labelSelector", k+"="+v)
	}

	res, err := req.Do().Get()

	if err != nil {
		return nil, err
	}

	return meta.ExtractList(res)
}

func pluralise(str string) string {

	exceptions := make(map[string]string)

	pluralise := namer.NewAllLowercasePluralNamer(exceptions)

	pluralType := types.Type{
		Name: types.Name{Name: str},
	}

	return pluralise.Name(&pluralType)
}

type objectRefSpec struct {
	Kind      string `yaml:"kind"`
	Namespace string `yaml:"namespace"`
	Labels    map[string]string
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
func (mf manifest) apply(client *kubernetes.Clientset, expectError bool) (*objectRef, error) {

	if valid := mf.valid(); !valid {
		return nil, fmt.Errorf("Invalid Manifest to apply")
	}

	var filePath string
	if testFilePath != "" {
		filePath = filepath.Dir(testFilePath) + "/" + mf.Path
	} else {
		filePath = mf.Path
	}

	ymlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err, "Error reading the "+mf.Path+" file.")
		return nil, err
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(ymlFile))
	if err != nil {
		log.Error(err, "Error reading test file.")
	}

	objRef := objectRef{
		Type: "Resource",
		Spec: objectRefSpec{
			Kind:      viper.GetString("kind"),
			Namespace: viper.GetString("metadata.namespace"),
			Labels:    viper.GetStringMapString("metadata.labels"),
		},
	}

	log.InfoWithFields(map[string]interface{}{
		"Action": mf.Action,
		"Path":   filePath,
	}, "Applying action to file")

	var cmd *exec.Cmd

	if mf.Action == "CREATE" {
		cmd = exec.Command("kubectl", "apply", "-f", filePath)
	} else {
		cmd = exec.Command("kubectl", "delete", "-f", filePath)
	}

	out, err := cmd.CombinedOutput()

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

	return &objRef, err
}
