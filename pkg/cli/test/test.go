// Package cmd test contains the logic for anchorctl test.
/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package test

import (
	"anchorctl/pkg/kubetest"
	"anchorctl/pkg/logging"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var description = `
test takes in the test steps as config and executes the tests in order.

Example usage: anchorctl test --file test.yaml --kind kubetest

Kinds of tests include:
- kubetest: Runs tests against the current context of a kube cluster. Requires --kubeconfig flag to be set.
    Types of tests include:
	- AssetJSONPath: Given jsonPath and value, assert that the jsonpath of object in cluster matches the value.
	- AssertValidation: Given action, filepath and expected error, apply the action to the object in file, and assert
		error is returned
	- AssertMutation: Given action, filepath, jsonPath and value, assert jsonpath after applying the object in the file.
`

var log = &logging.Logger{}

// Cmd represents the test command
var Cmd = &cobra.Command{
	Use:   "test",
	Short: "Command to run Anchor tests",
	Long:  description,
	Run:   test,
}

type TestKind string

var (
	KubeTest TestKind = "KubeTest"
)

func init() {
	var defaultKubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")

	// Local Flags
	Cmd.Flags().StringP("file", "f", "", "Input file with the tests.")
	Cmd.Flags().StringP("kubeconfig", "c", defaultKubeConfig, "Path to kubeconfig file.")
	Cmd.Flags().Float64P("threshold", "t", 80,
		"Percentage of tests to pass, else return failure.")
	Cmd.Flags().IntP("verbose", "v", 4,
		"Verbosity Level, choose between 1 being Fatal - 7 being .")
	Cmd.Flags().BoolP("incluster", "i", false, "Get kubeconfig from in cluster.")
}

func test(cmd *cobra.Command, args []string) {
	verbosity, err := cmd.Flags().GetInt("verbose")
	if err != nil {
		logrus.WithFields(logrus.Fields{"flag": "verbose"}).Error("Unable to parse flag. Defaulting to INFO.")
		verbosity = 4
	}
	log.SetVerbosity(verbosity)

	crd, testPath, err := decodeTestFile(cmd, log)
	if err != nil {
		log.Fatal(err, "Unable to decode test file.")
	}

	threshold, err := cmd.Flags().GetFloat64("threshold")
	if err != nil {
		log.Error(err, "Unable to parse flag. Defaulting to 100.")
		threshold = 100
	}

	test := getTest(crd, log, testPath, cmd)

	testResult := test.Assert()

	testResult.SetThreshold(threshold)
	testResult.Render()
}

func decodeTestFile(cmd *cobra.Command, log *logging.Logger) (*AnchorCRD, string, error) {
	testfile, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Fatal(err, "Failed to parse testfile.")
	}

	testContents, err := ioutil.ReadFile(filepath.Clean(testfile))
	if err != nil {
		log.Fatal(err, "Error reading the test file.")
	}

	testCRD := &AnchorCRD{}

	err = yaml.Unmarshal(testContents, &testCRD)
	if err != nil {
		return nil, "", err
	}

	return testCRD, testfile, nil
}

func getTest(crd *AnchorCRD, log *logging.Logger, testPath string, cmd *cobra.Command) AnchorTest {

	switch crd.Kind {

	case KubeTest:
		incluster, err := cmd.Flags().GetBool("incluster")
		if err != nil {
			log.Error(err, "Unable to parse flag. Defaulting to false.")
			incluster = false
		}

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		if err != nil && !incluster {
			log.Fatal(err, "Failed to parse kubeconfig flag")
		}

		kubeTest := &kubetest.KubeTest{
			Opts: kubetest.Options{
				Incluster:    incluster,
				Kubeconfig:   kubeconfig,
				TestFilepath: testPath,
				Logger:       log,
			},
		}

		err = mapstructure.Decode(crd.Spec, &kubeTest)
		if err != nil {
			log.Fatal(err, "Error parsing to KubeTest object")
		}

		log.Info("kind", "kubetest", "Decoded Tests")

		return kubeTest

	default:
		log.Fatal(fmt.Errorf("Unknown kind "+string(crd.Kind)), "Unknown AnchorTest kind. Accepted kinds: ")
		return nil
	}

}
