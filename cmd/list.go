package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/toddlers/l2info/lambdalayers"
	"github.com/toddlers/l2info/utils"
)

const layerInfo = `
{{range .}}-------------------------------------------------------------------------
Layer Name         : {{.LayerName}}
Description        : {{.LatestMatchingVersion.Description}}
Layer ARN          : {{.LayerArn}}
Layer Version ARN  : {{.LatestMatchingVersion.LayerVersionArn}}
Compatible Runtime : {{index .LatestMatchingVersion.CompatibleRuntimes 0}}
Created            : {{.LatestMatchingVersion.CreatedDate}}
{{end}}`

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list lambda layers",
	RunE:  listLayers,
}

func listLayers(cmd *cobra.Command, args []string) error {
	result, err := lambdalayers.ListLayers(runtime)
	if err != nil {
		utils.CheckError(err)
		return err
	}
	if len(result.Layers) == 0 {
		fmt.Println("No layer found with runtime : ", runtime)
		return nil
	}
	var report = template.Must(template.New("Layer Info").Parse(layerInfo))
	if err := report.Execute(os.Stdout, result.Layers); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&runtime, "runtime", "r", "python3.7", "compatible runtime")
}
