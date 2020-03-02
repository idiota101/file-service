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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete deletes the user",
	Long: `delete delets a file user.
	delete requires username and password to delete`,
	Run: func(cmd *cobra.Command, args []string) {
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return
		}

		resp, err := client.DeleteUser(requestCtx, &v1.DeleteUserRequest{
			Api: ApiVersion,
			UserDetails: &v1.User{
				Username: username,
				Password: password,
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
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("username", "u", "", "enter your username")
	deleteCmd.Flags().StringP("password", "p", "", "enter your password")
	deleteCmd.MarkFlagRequired("username")
	deleteCmd.MarkFlagRequired("password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
