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

	"github.com/anchorageio/anchorctl/utils/kubernetes"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("Could not parse kubeconfig flag", err.Error())
		}

		testfile, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("Could not parse testfile flag", err.Error())
		}

		client, err := kubernetes.GetKubeClient(false, kubeconfig)
		if err != nil {
			return fmt.Errorf("Could not get client", err.Error())
		}

		kubeTest, err := kubernetes.DecodeTestFile(testfile)
		if err != nil {
			return fmt.Errorf("Could not decode test file", err.Error())
		}

		pod, err := kubernetes.GetObject(client, &kubeTest.Object)
		if err != nil {
			return fmt.Errorf("Failed getting object", err)
		}

		for _, i := range kubeTest.Tests {
			switch i["type"]{

			case "AssertJSONPath":
				_, logs, err := kubernetes.AssertJsonpath(pod, i["jsonPath"], i["value"])

				if err != nil {
					return fmt.Errorf("AssertJsonPath Failed", err)
				}

				for _, i := range logs {
					cmd.Println(i)
				}

			default:
				cmd.Println(i["type"] + " is not a valid test type")

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
	testCmd.Flags().StringP("kubeconfig", "k", "~/.kube/config", "Path to kubeconfig file")
	testCmd.Flags().StringP("jsonpath", "j", "{.metadata.name}", "JSONPath string")
}
