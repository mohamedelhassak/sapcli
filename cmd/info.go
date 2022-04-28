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

var infoDesc string

func NewInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "info",
		Short:                 "Show info about SAPCLI tool",
		Long:                  `This command can be used view SAPCLI tool info & env`,
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			initInfoDesc()
			fmt.Println(infoDesc)

		},
	}

	return cmd
}
func initInfoDesc() {

	infoDesc = `
SAPCLI is a tool to trigger SAP Cloud APIs
	` + banner.Inline("have fun using sapcli") + SAPCLI_VERSION + `
	
### ENV VARIABLES
	SAPCLI_WORK_DIR :` + WORK_DIR + `
	
### FILES & DIRECTORIES
	Config file :` + viper.ConfigFileUsed() + `
	build  dir  :` + BUILDS_DIR + `
	logs dir    :` + LOGS_DIR
}
