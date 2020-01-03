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
	"bufio"
	"encoding/json"
	"fmt"
	"mqtt-ipfs/ipfs"
	"mqtt-ipfs/mqtt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var (
		clientId = "pibigstar"
		wg       sync.WaitGroup
	)
	client := mqtt.NewClient(clientId)
	err := client.Connect()
	if err != nil {
		//t.Errorf(err.Error())
		fmt.Println(err)
	}

	time.Sleep(time.Duration(1) * time.Second)

	go func() {
		err := client.Subscribe(func(c *mqtt.Client, msg *mqtt.Message) {
			fmt.Printf("Got message: %+v \n", msg)

			s := fmt.Sprintln(msg)

			hash := ipfs.UploadIPFS(s)
			fmt.Println("Hash for ", msg, "is :", hash)
			fmt.Println()

			wg.Done()
		}, 1, "mqtt")

		if err != nil {
			panic(err)
		}
	}()

	msg := &mqtt.Message{
		ClientID: "Client ID: " + clientId,
		Type:     "Message Type: " + "text",
		//Data:     "Hello Pibistar",
		//Time:     time.Now().UTC(),
	}

	/*for i := 0; i < 5; i++ {
		msg.Time = time.Now().UTC()
		data, _ := json.Marshal(msg)

		time.Sleep(time.Duration(1) * time.Second)
		wg.Add(1)
		err = client.Publish("mqtt", 1, false, data)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}*/

	for {
		time.Sleep(time.Duration(500) * time.Millisecond)
		fmt.Println("Please enter the message that you want to save on IPFS:")
		//fmt.Scanln(&msg.Data)

		//msg.Data = bufio.NewScanner(os.Stdin).Text()

		inputReader := bufio.NewReader(os.Stdin)
		inputData, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("There were errors reading, exiting program.")
			return
		}

		msg.Data = "Message: " + inputData

		msg.Time = "Timestamp: " + time.Now().UTC().String()
		data, _ := json.Marshal(msg)

		time.Sleep(time.Duration(1) * time.Second)
		wg.Add(1)
		err = client.Publish("mqtt", 1, false, data)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}


	time.Sleep(time.Duration(1) * time.Second)
	wg.Wait()
}
