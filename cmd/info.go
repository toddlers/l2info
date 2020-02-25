package cmd

import (
	"html/template"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/toddlers/l2info/lambdalayers"
	"github.com/toddlers/l2info/utils"
)

const layerInfoByVersion = `
Description        : {{ .Description}}
CompatibleRuntimes : {{index .CompatibleRuntimes 0}}
CodeSize           : {{ divide .Content.CodeSize }} {{ "KB"}}
CreatedDate        : {{ .CreatedDate}}
Layer Verion Arn   : {{ .LayerVersionArn}}
Location           : {{ .Content.Location}}
`

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get Information for the layer",
	RunE:  getInfo,
}

func getInfo(cmd *cobra.Command, args []string) error {
	fm := template.FuncMap{"divide": func(size int64) int64 {
		return size / 1000
	}}
	if layerArn != "" {
		layerInfo, err := lambdalayers.GetInfo(layerArn)
		if err != nil {
			utils.CheckError(err)
			return err
		}
		var report = template.Must(template.New("Layer Info By Version").Funcs(fm).Parse(layerInfoByVersion))
		if err := report.Execute(os.Stdout, layerInfo); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&layerArn, "layerArn", "l", "", "Layer Arn")
}
