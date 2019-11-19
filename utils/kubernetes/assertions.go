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

func assertDeny(){



}