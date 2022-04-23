/*
Copyright © 2022
This file is part of CLI application SAPCLI.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"moul.io/banner"
)

func NewInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "info",
		Short:                 "info",
		Long:                  `This command can be used view tool info & env`,
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("SAPCLI is a tool to trigger SAP Cloud APIs")
			fmt.Println(banner.Inline("have fun using sapcli"))
			fmt.Println("")
			fmt.Println("Version = " + SAPCLI_VERSION)
			fmt.Println("")
			fmt.Println("-->ENV VARIABLES")
			fmt.Println("SAPCLI_WORK_DIR =", WORK_DIR+"\n")
			fmt.Println("-->FILES & DIRECTORIES")
			fmt.Println("Config file :", viper.ConfigFileUsed())
			fmt.Println("build  dir  :", BUILDS_DIR)
			fmt.Println("logs   dir  :", LOGS_DIR)
			fmt.Println("")
			fmt.Println()
		},
	}

	return cmd
}
