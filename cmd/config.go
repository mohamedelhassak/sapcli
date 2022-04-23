/*
Copyright © 2022
This file is part of CLI application SAPCLI.
*/
package cmd

import (
	"errors"
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
var validConfigArgs = []string{"show"}

func isValidConfigArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Requires exactly 1 arg.")
	}
	return cobra.OnlyValidArgs(cmd, args)
}

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "config",
		Aliases:   []string{"conf"},
		Short:     "config",
		Long:      `This command can be used view config`,
		Example:   "",
		ValidArgs: validConfigArgs,
		Args:      isValidConfigArgs,
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
		Args:  cobra.NoArgs,

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
