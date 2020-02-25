package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/toddlers/l2info/lambdalayers"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the layer content",
	RunE:  downloadLayer,
}

func downloadLayer(cmd *cobra.Command, args []string) error {
	if exportPath == "" || layerArn == "" {
		fmt.Println("needed an export path for the content and layer arn")
	}
	layerDir, err := lambdalayers.DownloadLayer(exportPath, layerArn)
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded layer %s content to %s\n", layerArn, layerDir)

	return nil
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&exportPath, "exportPath", "e", exportPath, "Export path for the layer")
	downloadCmd.Flags().StringVarP(&layerArn, "layerArn", "l", "", "Lambda Layer ARN to download")
}
