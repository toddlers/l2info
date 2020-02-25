package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	runtime       string
	layerArn      string
	exportPath, _ = os.Getwd()
)

var rootCmd = &cobra.Command{
	Use:   "l2info",
	Short: "Lambda Layer Ops",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
