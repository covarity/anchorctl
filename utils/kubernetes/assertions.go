package kubernetes

import (
	"bytes"
	"k8s.io/client-go/util/jsonpath"
)

func AssertJsonpath(object interface{}, path, value string) (bool, []string, error) {

	jp := jsonpath.New("assertJsonpath")
	jp.AllowMissingKeys(true)
	err := jp.Parse("{" + path + "}")

	if err != nil {
		return false, nil, err
	}

	buf := new(bytes.Buffer)

	err = jp.Execute(buf, object)

	if err != nil {
		return false, nil, err
	}

	if buf.String() == value {
		return true, []string{"PASSED: " + path + " " + value}, nil
	}

	return false, []string{"FAILED: expected" + value + " got " + buf.String()}, nil

}
