/*
Copyright © 2022
This file is part of CLI application SAPCLI.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"moul.io/banner"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print version number of SAPCLI tool",
		Long:    `This command can be used get the version number of SAPCLI tool`,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(banner.Inline("have fun using sapcli"))
			fmt.Println(SAPCLI_VERSION)

		},
	}

	return cmd
}
