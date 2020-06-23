/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/nullck/k-logs-test/pkg/kubernetes_pods"
	"github.com/spf13/cobra"
)

var podName string
var namespaceName string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the test components",
	Long: `Execute run to start the test components. For example:

k-logs-test run --pod-name test-logs --namespace logs --elastic-endpoint https://localhost:9200/fluentd-2020`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("run called %v\n", podName)
		p, err := kubernetes_pods.CreatePod(podName, namespaceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(p)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&podName, "pod-name", "p", "", "The pod name")
	runCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "", "The pod namespace")
}
