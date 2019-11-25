package kubernetes

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/jsonpath"
)

func assertJsonpath(cmd *cobra.Command, object interface{}, path, value string) (bool, error) {

	jp := jsonpath.New("assertJsonpath")
	jp.AllowMissingKeys(true)
	err := jp.Parse("{" + path + "}")
	passed := true

	if err != nil {
		cmd.PrintErrln("Cannot parse jsonpath. ", err)
		return false, err
	}

	buf := new(bytes.Buffer)

	objects := getSlice(cmd, object)

	for _, i := range objects {

		err = jp.Execute(buf, i)
		if err != nil {
			cmd.PrintErrln("Error executing jsonpath on object. ", err)
			passed = false
			break
		}

		if buf.String() == value {
			cmd.Println("PASSED: " + path + " " + value)
		} else {
			cmd.Println("FAILED: expected" + value + " got " + buf.String())
			passed = false
		}

		buf.Reset()

	}

	return passed, err
}

func assertValidation(client *kubernetes.Clientset, action, filepath, expectedError string) bool {
	_, _, err := applyAction(client, filepath, action)

	if err != nil && err.Error() == expectedError {
		fmt.Println("Passed validation")
		return true
	}

	fmt.Println("Failed Validation, got no error, expected error = ", expectedError)
	return false

}

func assertMutation(client *kubernetes.Clientset, cmd *cobra.Command, action, filepath, jsonpath, value string) (bool, error) {
	_, obj, err := applyAction(client, filepath, action)
	if err != nil {
		return false, fmt.Errorf("Errored out: ", err)
	}

	return assertJsonpath(cmd, obj, jsonpath, value)

}

func assertNetworkPolicies(sourceNamespace, sourceLabelKey, sourceLabelValue, destNamespace, destNamespaceKey, destValue, port, ipaddress string){
// TODO: If ip address is provided, check it returns a 200 with the port
// TODO: Else, create pod in sourceNamespace with sourceLabelKey and sourceLabelValue
// TODO: Create pod in destNamespace with destLabelKey and destLabelValue
// TODO: Exec into source pod, telnet destination pod in the provided IP address
}
