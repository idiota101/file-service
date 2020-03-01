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
	"google.golang.org/grpc/metadata"

	v1 "github.com/sajanjswl/file-service/pkg/api/v1"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload uploads a file to the mongoDB database",
	Long: `upload a documents on the server through gRPC.
	upload requires a specific file path`,
	Run: func(cmd *cobra.Command, args []string) {

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return
		}

		file, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		ctx := metadata.AppendToOutgoingContext(requestCtx, "username", username, "password", password)

		stream, err := client.UploadFile(ctx)
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 2048)

		for {

			n, err := file.Read(buf)

			if err != nil {
				if err == io.EOF {

					break
				}
				log.Fatal(err)
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
	uploadCmd.Flags().StringP("username", "u", "", "username")

	uploadCmd.Flags().StringP("password", "p", "", "password")

	uploadCmd.Flags().StringP("path", "e", "", "file path")
	registerCmd.MarkFlagRequired("path")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
