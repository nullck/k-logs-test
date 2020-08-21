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
	"log"
	"time"

	"github.com/nullck/k-logs-test/pkg/elastic"
	"github.com/nullck/k-logs-test/pkg/kubernetes_pods"
	"github.com/nullck/k-logs-test/pkg/slack"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the test components",
	Long: `Execute run to start the test components. For example:

k-logs-test run --pod-name test-logs --logs-hits 30 --namespace logs --elastic-endpoint https://localhost:9200/fluentd-2020 --slack-alert-enabled true --threshold 10 --webhook-url https://hooks.slack.com/services/XXX --channel #general`,
	Run: func(cmd *cobra.Command, args []string) {
		var podName, namespaceName, elasticAddr, elasticRes, slackChannel, slackWebhookUrl, slackMsg string
		var logsHits, threshold int
		var slackAlertEnabled bool

		type s = slack.Slack
		type p = kubernetes_pods.Pod
		type e = elastic.ES

		var (
			po = p{
				PodName:       podName,
				NamespaceName: namespaceName,
			}

			es = e{
				ElasticAddr: elasticAddr,
				PodName:     podName,
				LogsHits:    logsHits,
				Threshold:   threshold,
			}

			sl = s{
				WebhookUrl: slackWebhookUrl,
				Username:   "k-logs",
				Channel:    slackChannel,
			}
		)

		_, err := po.CreatePod(logsHits)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("k-logs checking total pods logs %d ...\n", logsHits)
		time.Sleep(time.Duration(logsHits) * time.Second)

		elasticRes, err = es.Search()
		log.Printf("status: %v\n", elasticRes)

		if elasticRes == "ALERT" {
			if slackAlertEnabled {
				slackMsg = "error: k-logs threshold reached!"
				err = sl.Notification(slackMsg)
				log.Printf("slack notification sent: %v\n", slackMsg)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}

		_, err = po.DeletePod()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&podName, "pod-name", "p", "k-logs-test", "The pod name")
	runCmd.Flags().IntVar(&logsHits, "logs-hits", 30, "The number of logs hits")
	runCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "default", "The pod namespace")
	runCmd.Flags().StringVarP(&elasticAddr, "elastic-endpoint", "e", "https://localhost:9200/fluentd", "The ElasticSearch Endpoint and the logs index name")
	runCmd.Flags().BoolVarP(&slackAlertEnabled, "slack-alert-enabled", "a", false, "Enable or not slack alerts")
	runCmd.Flags().IntVar(&threshold, "threshold", 1000, "The Alert Threshould in milliseconds")
	runCmd.Flags().StringVarP(&slackChannel, "channel", "c", "#k-logs", "The Slack Channel for notification")
	runCmd.Flags().StringVarP(&slackWebhookUrl, "webhook-url", "w", "", "The Slack Webhook Url for notification")
}
