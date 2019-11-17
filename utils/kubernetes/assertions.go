package kubernetes

import (
	"bytes"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/jsonpath"
)

func assertJsonpath(cmd *cobra.Command, object interface{}, path, value string) (bool, error) {

	jp := jsonpath.New("assertJsonpath")
	jp.AllowMissingKeys(true)
	err := jp.Parse("{" + path + "}")

	if err != nil {
		cmd.PrintErrln("Cannot parse jsonpath. ", err)
		return false, err
	}

	buf := new(bytes.Buffer)

	err = jp.Execute(buf, object)

	if err != nil {
		cmd.PrintErrln("Error executing jsonpath on object. ", err)
		return false, err
	}

	if buf.String() == value {
		cmd.Println("PASSED: " + path + " " + value)
		return true, nil
	}

	cmd.Println("FAILED: expected" + value + " got " + buf.String())

	return false, nil

}
