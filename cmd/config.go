package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/mohamedelhassak/sapcli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// config model
type Config struct {
	Creds struct {
		APIToken       string `mapstructure:"api-token"`
		SubscriptionID string `mapstructure:"subscription-id"`
	} `mapstructure:"creds"`
}

//init config struct
var cfg Config

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf"},
		Short:   "config",
		Long:    `This command can be used view config`,
		Example: "",
		Args:    cobra.ExactArgs(1),
	}
	cmd.AddCommand(
		NewConfigShowCmd(),
	)
	return cmd
}

func NewConfigShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show",
		Long:  `This command can be used to show subscription & api token`,
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			fileName := viper.ConfigFileUsed()
			extension := filepath.Ext(fileName)

			err := viper.Unmarshal(&cfg)
			if err != nil {
				return
			}

			fmt.Println(utils.PrettyPrint(cfg, extension))

		},
	}
	return cmd
}
