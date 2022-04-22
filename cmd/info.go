package cmd

import (
	"fmt"

	"github.com/mohamedelhassak/sapcli/utils"
	"github.com/spf13/cobra"
	"moul.io/banner"
)

func NewInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "info",
		Long:  `This command can be used view application info & env`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Println(utils.UnknownCommandMsg("info"))
				return
			}

			fmt.Println("SAPCLI is a tool to trigger SAP Cloud APIs")
			fmt.Println(banner.Inline("have fun using sapcli"))
			fmt.Println("")
			fmt.Println("Version = " + SAPCLI_VERSION)
			fmt.Println("")
			fmt.Println("-->ENV VARIABLES")
			fmt.Println("SAPCLI_WORK_DIR =", WORK_DIR+"\n")
			fmt.Println("-->FILES & DIRECTORIES")
			fmt.Println("Config file :", CONF_FILE_NAME)
			fmt.Println("build  dir  :", BUILDS_DIR)
			fmt.Println("logs   dir  :", LOGS_DIR)
			fmt.Println("")
			fmt.Println("")
		},
	}

	return cmd
}
