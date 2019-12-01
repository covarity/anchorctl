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
package cmd

import (
	"github.com/anchorageio/anchorctl/utils/kubernetes"
	"github.com/anchorageio/anchorctl/utils/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var description = `
Test takes in the test steps as config and executes the tests in order.

Example usage: anchorctl test --file test.yaml --kind kubetest

Kinds of tests include:
- kubetest: Runs tests against the current context of a kube cluster. Requires --kubeconfig flag to be set.
    Types of tests include:
	- AssetJSONPath: Given jsonPath and value, assert that the jsonpath of object in cluster matches the value.
	- AssertValidation: Given action, filepath and expected error, apply the action to the object in file, and assert error is returned
	- AssertMutation: Given action, filepath, jsonPath and value, assert jsonpath after applying the object in the file.
`

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Command to run Anchor tests",
	Long:  description,
	Run:   testExecute,
}

func init() {
	rootCmd.AddCommand(testCmd)
	var defaultKubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	// Local Flags
	testCmd.Flags().StringP("file", "f", "", "Input file with the tests.")
	testCmd.Flags().StringP("kubeconfig", "c", defaultKubeConfig, "Path to kubeconfig file.")
	testCmd.Flags().StringP("kind", "k", "kubetest", "Kind of test, only kubetest is supported at the moment.")
	testCmd.Flags().Float64P("threshold", "t", 80, "Percentage of tests to pass, else return failure.")
	testCmd.Flags().IntP("verbose", "v", 5, "Verbosity Level, choose between 1 being Fatal - 7 being .")
	testCmd.Flags().BoolP("incluster", "i", false, "Get kubeconfig from in cluster.")
}

func testExecute(cmd *cobra.Command, args []string) {
	verbosity, err := cmd.Flags().GetInt("verbose")
	if err != nil {
		logrus.WithFields(logrus.Fields{"flag": "verbose"}).Error("Unable to parse flag. Defaulting to INFO.")
		verbosity = 5
	}
	log := &logging.Logger{}
	log.SetVerbosity(verbosity)

	kind, err := cmd.Flags().GetString("kind")
	if err != nil {
		log.Error(err, "Unable to parse flag. Defaulting to Kubetest.")
		kind = "kubetest"
	}

	testfile, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Fatal(err, "Failed to parse testfile.")
	}

	threshold, err := cmd.Flags().GetFloat64("threshold")
	if err != nil {
		log.Error(err, "Unable to parse flag. Defaulting to 100.")
		threshold = 100
	}

	incluster, err := cmd.Flags().GetBool("incluster")
	if err != nil {
		log.Error(err, "Unable to parse flag. Defaulting to false.")
		incluster = false
	}

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil && incluster == false {
		log.Fatal(err, "Failed to parse kubeconfig flag")
	}

	switch kind {

	case "kubetest":
		log.Info("kind", "kubetest", "Starting Tests")
		kubernetes.Assert(log, threshold, incluster, kubeconfig, testfile)
		if err != nil {
			log.Fatal(err, "Failed Tests")
		}
		log.Info("kind", "kubetest", "Finished Tests")
	}
}
