package resource

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

type Manifest struct {
	Path   string `yaml:"path"`
	Action string `yaml:"action"`
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
