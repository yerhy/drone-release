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
	im "drone-release/tools/images"

	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "image command for nscloud images",
}

var genConfigCmd = &cobra.Command{
	Use:   "generate-svnc-config",
	Short: "generate image-sync config.json",
	Run: func(cmd *cobra.Command, args []string) {
		im.GenerateSyncConfig()
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(genConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genConfigCmd.Flags().StringVarP(&im.SyncRecordFile, "record", "r", "svn_record.toml", "svn codes record file")
}
