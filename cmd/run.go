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
	"log"
	"strings"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/nullck/k-logs-test/pkg/elastic"
	"github.com/nullck/k-logs-test/pkg/kubernetes_pods"
	"github.com/nullck/k-logs-test/pkg/slack"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename = "k-logs"
	envPrefix             = "K_LOGS"
)

func NewRootCommand() *cobra.Command {
	var namespaceName, elasticAddr, elasticRes, slackChannel, slackWebhookUrl, slackMsg, promGWAddr string
	var logsHits, threshold, promGWPort int
	var slackAlertEnabled, promEnabled bool

	// runCmd represents the run command
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start the test components",
		Long: `Execute run to start the test components. For example:
		k-logs-test run --logs-hits 30 --namespace logs --elastic-endpoint https://localhost:9200/fluentd-2020 --prom true --prom-addr prometheus-pushgateway --prom-port 9091 --slack-alert true --threshold 10 --webhook-url https://hooks.slack.com/services/XXX --channel #general`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			podName := fmt.Sprintf("k-logs-%s", GeneratePodName())

			var (
				po = kubernetes_pods.Pod{
					PodName:       podName,
					NamespaceName: namespaceName,
				}

				es = elastic.ES{
					ElasticAddr: elasticAddr,
					PodName:     podName,
					LogsHits:    logsHits,
					Threshold:   threshold,
				}

				sl = slack.Slack{
					WebhookUrl: slackWebhookUrl,
					Username:   "k-logs",
					Channel:    slackChannel,
				}
			)
			pop := &po
			esp := &es
			_, err := pop.CreatePod(logsHits)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("k-logs checking total pods logs %d ...\n", logsHits)
			time.Sleep(time.Duration(logsHits) * time.Second)

			elasticRes, err = esp.Search(promEnabled, promGWAddr, promGWPort)
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

			_, err = pop.DeletePod(pop.PodName)
			if err != nil {
				log.Fatalln(err)
			}
			// cleaner process, in case we are still some missing pods
			pop.Cleaner()
		},
	}
	runCmd.Flags().IntVar(&logsHits, "logs-hits", 30, "The number of logs hits")
	runCmd.Flags().StringVarP(&namespaceName, "namespace", "n", "default", "The pod namespace")
	runCmd.Flags().StringVarP(&elasticAddr, "elastic-endpoint", "e", "https://localhost:9200/fluentd", "ElasticSearch Endpoint and the logs index name")
	runCmd.Flags().BoolVar(&promEnabled, "prom-enabled", false, "Enable or not the prometheus metrics")
	runCmd.Flags().StringVar(&promGWAddr, "prom-endpoint", "prometheus-pushgateway", "The prometheus gateway addr")
	runCmd.Flags().IntVar(&promGWPort, "prom-port", 9091, "The prometheus gateway port")
	runCmd.Flags().BoolVarP(&slackAlertEnabled, "slack-alert-enabled", "a", false, "Enable or not slack alerts")
	runCmd.Flags().IntVar(&threshold, "threshold", 1000, "The Alert Threshould in milliseconds")
	runCmd.Flags().StringVarP(&slackChannel, "channel", "c", "#k-logs", "The Slack Channel for notification")
	runCmd.Flags().StringVarP(&slackWebhookUrl, "webhook-url", "w", "", "The Slack Webhook Url for notification")

	return runCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscore starting by KLOGS_
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func GeneratePodName() string {
	petname.NonDeterministicMode()
	podName := petname.Generate(2, "-")
	return podName
}
