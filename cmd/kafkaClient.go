/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var (
	command    string
	tlsEnable  bool
	hosts      string
	topic      string
	partition  int
	clientcert string
	clientkey  string
	cacert     string
)

// kafkaClientCmd represents the kafkaClient command
var kafkaClientCmd = &cobra.Command{
	Use:   "kafkaClient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		kafkaA, err := cmd.Flags().GetString("kafkaA")
		if err != nil {
			return errors.New("please input addr name ")
		}
		if len(kafkaA) == 0 {
			return errors.New("请携带参数-k 或者 --kafkaA")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		kafkaA, _ := cmd.Flags().GetString("kafkaA")
		flag.StringVar(&command, "command", "consumer", "consumer|producer")
		flag.BoolVar(&tlsEnable, "tls", false, "TLS enable")
		flag.StringVar(&hosts, "host", kafkaA, "Common separated kafka hosts")
		flag.StringVar(&topic, "topic", "JCJ_HUAWU", "Kafka topic")
		flag.IntVar(&partition, "partition", 0, "Kafka topic partition")
		flag.StringVar(&clientcert, "cert", "cert.pem", "Client Certificate")
		flag.StringVar(&clientkey, "key", "key.pem", "Client Key")
		flag.StringVar(&cacert, "ca", "ca.pem", "CA Certificate")
		flag.Parse()

		config := sarama.NewConfig()
		if tlsEnable {
			//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
			tlsConfig, err := genTLSConfig(clientcert, clientkey, cacert)
			if err != nil {
				log.Fatal(err)
			}

			config.Net.TLS.Enable = false
			config.Net.TLS.Config = tlsConfig
		}
		client, err := sarama.NewClient(strings.Split(hosts, ","), config)
		if err != nil {
			log.Fatalf("unable to create kafka client: %q", err)
		}

		if command == "consumer" {
			consumer, err := sarama.NewConsumerFromClient(client)
			if err != nil {
				log.Fatal(err)
			}
			defer consumer.Close()
			loopConsumer(consumer, topic, partition)
		} else {
			producer, err := sarama.NewAsyncProducerFromClient(client)
			if err != nil {
				log.Fatal(err)
			}
			defer producer.Close()
			loopProducer(producer, topic, partition)
		}
		fmt.Println("kafkaClient called")
	},
}

func genTLSConfig(clientcertfile, clientkeyfile, cacertfile string) (*tls.Config, error) {
	// load client cert
	clientcert, err := tls.LoadX509KeyPair(clientcertfile, clientkeyfile)
	if err != nil {
		return nil, err
	}

	// load ca cert pool
	cacert, err := ioutil.ReadFile(cacertfile)
	if err != nil {
		return nil, err
	}
	cacertpool := x509.NewCertPool()
	cacertpool.AppendCertsFromPEM(cacert)

	// generate tlcconfig
	tlsConfig := tls.Config{}
	tlsConfig.RootCAs = cacertpool
	tlsConfig.Certificates = []tls.Certificate{clientcert}
	tlsConfig.BuildNameToCertificate()
	// tlsConfig.InsecureSkipVerify = true // This can be used on test server if domain does not match cert:
	return &tlsConfig, err
}

func loopProducer(producer sarama.AsyncProducer, topic string, partition int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
		} else if text == "exit" || text == "quit" {
			break
		} else {
			producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(text)}
			log.Printf("Produced message: [%s]\n", text)
		}
		fmt.Print("> ")
	}
}

func loopConsumer(consumer sarama.Consumer, topic string, partition int) {
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	if err != nil {
		log.Println(err)
		return
	}
	defer partitionConsumer.Close()

	for {
		msg := <-partitionConsumer.Messages()
		log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)
	}
}

func init() {
	kafkaClientCmd.PersistentFlags().StringP("kafkaA", "k", "", "input kafka addr")
	rootCmd.AddCommand(kafkaClientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kafkaClientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kafkaClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
