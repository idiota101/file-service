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
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "downloads the file saved in mongoDB",
	Long: `downloads the file saved in mongoDB
	downloads requires username and password  `,
	Run: func(cmd *cobra.Command, args []string) {

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return
		}

		fmt.Println("starting to do a server Streaming  RPC...")

		req := &v1.DownloadFileRequest{
			Api:      ApiVersion,
			Username: username,
			Password: password,
		}

		resStream, err := client.DownloadFile(requestCtx, req)

		if err != nil {
			log.Fatalf("error while calling DownloadFile Rpc %v", err)
		}

		file, err := os.Create(string(username + ".pdf"))

		if err != nil {
			log.Fatal("err")
		}
		for {
			b, err := resStream.Recv()

			if err != nil {
				if err == io.EOF {
					log.Println("successfully downloaded")
					break
				}
				log.Fatal("unable to download file")

			}

			_, err = file.Write(b.GetContent())

			if err != nil {
				log.Fatalln(err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("username", "u", "", "enter your username")
	downloadCmd.Flags().StringP("password", "p", "", "enter your password")
	downloadCmd.MarkFlagRequired("username")
	downloadCmd.MarkFlagRequired("password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
