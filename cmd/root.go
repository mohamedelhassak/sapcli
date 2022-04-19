package cmd

import (
	"github.com/spf13/cobra"
)

var config string
var rootCmd = &cobra.Command{
	Use:   "sapcli",
	Short: "sapcli: A sample CLI to build and deploy CX Hybris to SAP Cloud",
	Long:  `sapcli: A sample CLI to build and deploy CX Hybris to SAP Cloud`,
}

func init() {

	rootCmd.PersistentFlags().String("config", "", "To set config file, default is $SAPCLI_WORK_DIR/.config.(yaml/json)")
}

// Execute executes the root command.
func Execute() error {

	return rootCmd.Execute()
}
