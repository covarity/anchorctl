package kubernetes

import (
	"bytes"
	"fmt"
	"github.com/anchorageio/anchorctl/utils/logging"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/jsonpath"
)

func assertJsonpath(logging *logging.Logger, object interface{}, path, value string) (bool, error) {

	jp := jsonpath.New("assertJsonpath")
	jp.AllowMissingKeys(true)
	err := jp.Parse("{" + path + "}")
	passed := true
	logger := logging.Log

	if err != nil {
		logger.WithFields(log.Fields{ "jsonPath": path,}).Error("Cannot parse JSONPath")
		return false, err
	}

	buf := new(bytes.Buffer)

	objects := getSlice(object)

	for _, i := range objects {

		err = jp.Execute(buf, i)
		if err != nil || buf.String() != value {
			logger.WithFields(log.Fields{ "jsonPath": path,}).Error("Cannot execute JSONPath")
			passed = false
			break
		}
		buf.Reset()

	}

	return passed, err
}

func assertValidation(client *kubernetes.Clientset, action, filepath, expectedError string) bool {
	_, _, err := applyAction(client, filepath, action)

	if err != nil && err.Error() == expectedError {
		return true
	}

	return false

}

func assertMutation(client *kubernetes.Clientset, logging *logging.Logger, action, filepath, jsonpath, value string) (bool, error) {
	_, obj, err := applyAction(client, filepath, action)
	if err != nil {
		return false, fmt.Errorf("Errored out: ", err)
	}

	return assertJsonpath(logging, obj, jsonpath, value)

}

func assertNetworkPolicies(sourceNamespace, sourceLabelKey, sourceLabelValue, destNamespace, destNamespaceKey, destValue, port, ipaddress string){
	// TODO: If ip address is provided, check it returns a 200 with the port
	// TODO: Else, create pod in sourceNamespace with sourceLabelKey and sourceLabelValue
	// TODO: Create pod in destNamespace with destLabelKey and destLabelValue
	// TODO: Exec into source pod, telnet destination pod in the provided IP address
}
