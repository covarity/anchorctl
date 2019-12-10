package kubetest

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/jsonpath"
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

func decodeTestFile(client *kubernetes.Clientset, filePath string) (*kubeTest, error) {
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

func assertJSONPath(objs []runtime.Object, path, value string) (bool, error) {

	jp := jsonpath.New("assertJsonpath")
	jp.AllowMissingKeys(true)
	err := jp.Parse("{" + path + "}")
	passed := true
	if err != nil {
		log.Error(err, "Cannot parse JSONPath")
		return false, err
	}

	buf := new(bytes.Buffer)

	for _, i := range objs {
		err = jp.Execute(buf, i)
		if err != nil {
			log.Error(err, "Cannot execute JSONPath")
			passed = false
			break
		} else if buf.String() != value {
			log.WarnWithFields(map[string]interface{}{
				"jsonpath": path,
				"expected": value,
				"got":      buf.String(),
				"status":   "FAILED",
			}, "Failed asserting jsonpath on obj")
			passed = false
			break
		}
		buf.Reset()
	}
	log.InfoWithFields(map[string]interface{}{
		"test":   "AssertJSONPath",
		"path":   path,
		"status": "PASSED",
	}, "JSON Path matches expected value.")

	return passed, err
}
