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

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload uploads a file to the mongoDB database",
	Long: `upload a documents on the server through gRPC.
	upload requires a specific file path`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return
		}

		file, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		stream, err := client.UploadFile(requestCtx)

		buf := make([]byte, 1024)

		for {

			n, err := file.Read(buf)

			if err != nil {
				if err == io.EOF {

					break
				}
				log.Println(err)
			}

			log.Println("Sendin", n, "bytes", "...")
			err = stream.Send(&v1.Chunk{

				Content: buf[:n],
			})

			if err != nil {
				if err == io.EOF {

					break
				}

				log.Fatalf("Send(%v) = %v", stream, err)
			}

		}

		reply, err := stream.CloseAndRecv()

		fmt.Println(reply)

	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringP("path", "p", "", "file path")
	registerCmd.MarkFlagRequired("path")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
