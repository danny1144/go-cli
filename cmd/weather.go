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
	"errors"
	"fmt"
	"my-client/server"
	"my-client/tools"
	"os"

	"github.com/spf13/cobra"
)

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		code, err := cmd.Flags().GetInt("code")
		if err != nil {
			return err
		}
		if len(name) == 0 && code == 0 {
			return errors.New("请携带参数-n name 或者 -c code")
		}
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("weather called")
		code, _ := cmd.Flags().GetInt("code")

		if code > 0 {
			info, err := server.GetWeather(code)
			if err != nil {
				errors.New(err.Error())
				os.Exit(1)
			}
			fmt.Println("current city weather is ", info)
			return
		}
		name, _ := cmd.Flags().GetString("name")
		if name != "" {
			code := tools.CityMap[name]
			info, err := server.GetWeather(code)
			if err != nil {
				errors.New(err.Error())
				os.Exit(1)
			}
			fmt.Println("current city weather is ", info)
		}
	},
}

func init() {
	weatherCmd.PersistentFlags().StringP("name", "n", "", "please input city name")
	weatherCmd.PersistentFlags().IntP("code", "c", 0, "please input city code ")
	rootCmd.AddCommand(weatherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// weatherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weatherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
