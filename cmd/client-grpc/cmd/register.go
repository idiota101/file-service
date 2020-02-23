// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
	"github.com/spf13/cobra"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register is used to register a file user",
	Long: `register registers a file user.
	register requires username emailId and password to register`,
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := client.RegisterUser(requestCtx, &v1.RegisterUserRequest{
			Api: apiVersion,
			UserDetails: &v1.UserDetails{
				Name:     "sajan",
				Email:    "sjnjaiswalfhgfjfg",
				Password: "sjfddnf",
			},
		})

		if err != nil {
			log.Println(err)
		} else {

			log.Println(resp)
		}

	},
}

func init() {
	RootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringP("name", "n", "", "enter your name")
	registerCmd.Flags().StringP("emailID", "u", "", "enter your email address")
	registerCmd.Flags().StringP("password", "p", "", "enter your password")
	registerCmd.MarkFlagRequired("emailID")
	registerCmd.MarkFlagRequired("password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
