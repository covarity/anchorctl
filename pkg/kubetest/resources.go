package kubetest

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"os/exec"
	"path/filepath"
)

func (ob objectRef) getObject(client *kubernetes.Clientset) ([]runtime.Object, error) {

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
func (mf manifest) apply(expectError bool) (*objectRef, error) {

	if valid := mf.valid(); !valid {
		return nil, fmt.Errorf("invalid Manifest to apply")
	}

	var filePath string
	if testFilePath != "" {
		filePath = filepath.Clean(filepath.Join(filepath.Dir(testFilePath), mf.Path))
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

func (mf manifest) getObjectref(filePath string) (*objectRef, error) {
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
