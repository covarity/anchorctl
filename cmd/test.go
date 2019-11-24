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
	"os"
	"github.com/sirupsen/logrus"
	"github.com/anchorageio/anchorctl/utils/kubernetes"
	"github.com/anchorageio/anchorctl/utils/logging"
	"github.com/spf13/cobra"
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
	Long: description,
	Run: testExecute,
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Local Flags
	testCmd.Flags().StringP("file", "f", "", "Input file with the tests.")
	testCmd.Flags().StringP("kubeconfig", "c", "~/.kube/config", "Path to kubeconfig file.")
	testCmd.Flags().StringP("kind", "k", "kubetest", "Kind of test, only kubetest is supported at the moment.")
	testCmd.Flags().Float64P("threshold", "t", 80, "Percentage of tests to pass, else return failure.")
	testCmd.Flags().IntP("verbose", "v", 3, "Verbosity Level, choose between 1 being Fatal - 7 being .")
}

func testExecute(cmd *cobra.Command, args []string) {
	verbosity, err := cmd.Flags().GetInt("verbose")
	if err != nil {
		logrus.WithFields(logrus.Fields{ "flag": "verbose"}).Error("Unable to parse flag. Defaulting to INFO.")
		verbosity = 5
	}
	log := &logging.Logger{}
	log.SetVerbosity(verbosity)
	logger := log.GetLogger()

	kind, err := cmd.Flags().GetString("kind")
	if err != nil {
		logger.WithFields(logrus.Fields{ "flag": "kind"}).Error("Unable to parse flag. Defaulting to kubetest.")
		kind = "kubetest"
	}

	testfile, err := cmd.Flags().GetString("file")
	if err != nil {
		logger.WithFields(logrus.Fields{ "flag": "file"}).Fatal("Unable to parse flag.")
	}

	threshold, err := cmd.Flags().GetFloat64("threshold")
	if err != nil {
		logger.WithFields(logrus.Fields{ "flag": "threshold"}).Error("Unable to parse flag. Defaulting to 100.")
	}

	switch kind {

	case "kubetest":
		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		if err != nil {
			logger.WithFields(logrus.Fields{ "flag": "kubeconfig"}).Fatal("Unable to parse flag.")
		}

		logger.WithFields(logrus.Fields{ "kind": "kubetest"}).Info("Starting kube assert")
		err = kubernetes.Assert(log, threshold, kubeconfig, testfile)
		if err != nil {
			os.Exit(1)
		}

	}
}
