package cmd

import (
	"fmt"
	"log"

	"github.com/mohamedelhassak/sapcli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "sapcli",
	Short: "sapcli: A sample CLI to build and deploy CX Hybris to SAP Cloud",
	Long:  `sapcli: A sample CLI to build and deploy CX Hybris to SAP Cloud`,
}

func init() {
	cobra.OnInitialize(initConfigs)
	//rootCmd.PersistentFlags().String("config", "")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "To set config file, default is $SAPCLI_WORK_DIR/.config.yaml")
}

// Execute executes the root command.
func Execute() error {

	return rootCmd.Execute()
}

//set custom config file if --config flag passed in cli
func initConfigs() {

	//if --config passed
	if cfgFile != "" {

		//check if file exist
		if !utils.IsFileOrDirExists(cfgFile) {
			log.Fatalf("[ERROR!...] file not found:  %s", cfgFile)
		}

		viper.SetConfigFile(cfgFile)
		fmt.Println("[INFO!...] Using config file : " + viper.ConfigFileUsed())

	} else {

		CONF_FILE_NAME = utils.SearchFileByPattern(CONF_FILE_NAME_PATTERN, WORK_DIR)
		//check if file exist
		if !utils.IsFileOrDirExists(CONF_FILE_NAME) {
			log.Fatalf("[ERROR!...] file not found:  %s", CONF_FILE_NAME)
		}

		viper.SetConfigFile(CONF_FILE_NAME)
		fmt.Println("[INFO!...] Using default config file : " + viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("[FAILED!...] Failed to read de config file: %v", err)
	}

	SAP_CLOUD_API_URL = "https://portalrotapi.hana.ondemand.com/v2/subscriptions/" + viper.Get("creds.subscription-id").(string)
	API_TOKEN = viper.Get("creds.api-token").(string)

}
