package kubernetes

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Assert(cmd *cobra.Command, kubeconfig, testfile string) error {
	client, err := getKubeClient(false, kubeconfig)

	if err != nil {
		return fmt.Errorf("Could not get client", err.Error())
	}

	kubeTest, err := decodeTestFile(testfile)
	if err != nil {
		return fmt.Errorf("Could not decode test file", err.Error())
	}

	object, err := getObject(client, kubeTest)
	if err != nil {
		return fmt.Errorf("Failed getting object", err)
	}

	for _, i := range kubeTest.Tests {
		switch i["type"]{

		case "AssertJSONPath":
			_, err := assertJsonpath(cmd, object, i["jsonPath"], i["value"])

			if err != nil {
				return fmt.Errorf("AssertJsonPath Failed", err)
			}

		case "AssertValidation":

		default:
			cmd.Println(i["type"] + " is not a valid test type")

		}
	}

	return nil
}
