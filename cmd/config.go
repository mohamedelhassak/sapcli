package cmd

import (
	"fmt"

	"github.com/mohamedelhassak/sapcli/utils"
	"github.com/spf13/cobra"
)

// config model
type Config struct {
	Creds struct {
		ApiToken       string `yaml:"api-token"`
		SubscriptionId string `yaml:"subscription-id"`
	} `yaml:"creds"`
}

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "config",
	Long:    `This command can be used view config`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show",
	Long:  `This command can be used to show subscription & api token`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			fmt.Println(utils.UnknownCommandMsg("config show"))
			return
		}
		readYAMLConfigs(CONF_FILE_NAME, &cfg)
		fmt.Println("[INFO!...] Using defaul config file :" + CONF_FILE_NAME)
		fmt.Println(utils.PrettyPrintYAML(cfg))

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
}
