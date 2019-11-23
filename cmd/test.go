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
	"fmt"
	"os"

	"github.com/anchorageio/anchorctl/utils/kubernetes"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Command to run Anchor tests",
	Long: `Test takes in the test steps as config and executes the tests in order.

Example usage: anchorctl test -f test.yaml

Kinds of tests include:
- kubetest: Requires --kubeconfig flag to be set. Runs tests against the current context of a kube cluster.
			Types of tests include: 
			- AssetJSONPath: Example fields are type: AssertJSONPath, jsonPath: ".spec.nodeName", value: "docker-desktop"`,
	RunE: func(cmd *cobra.Command, args []string) error {

		kind, err := cmd.Flags().GetString("kind")
		if err != nil {
			return fmt.Errorf("Could not parse kind flag", err.Error())
		}

		testfile, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("Could not parse testfile flag", err.Error())
		}

		threshold, err := cmd.Flags().GetFloat64("threshold")
		if err != nil {
			return fmt.Errorf("Could not parse threshold flag", err.Error())
		}

		switch kind {

		case "kubetest":
			kubeconfig, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return fmt.Errorf("Could not parse kubeconfig flag", err.Error())
			}

			err = kubernetes.Assert(cmd, threshold, kubeconfig, testfile)
			if err != nil {
				os.Exit(1)
			}

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	testCmd.Flags().StringP("file", "f", "", "Input file with the tests.")
	testCmd.Flags().StringP("kubeconfig", "c", "~/.kube/config", "Path to kubeconfig file.")
	testCmd.Flags().StringP("kind", "k", "kubetest", "Kind of test, only kubetest is supported at the moment.")
	testCmd.Flags().Float64P("threshold", "t", 80, "Percentage of tests to pass, else return failure.")
	testCmd.Flags().IntP("verbose", "v", 3, "Verbosity Level, choose between 1 being Fatal - 7 being .")
}
