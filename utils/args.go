package utils

import (
	"errors"

	"github.com/spf13/cobra"
)

//return error if more than one arg passed
func IsOneAndOnlyValidArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("\"" + cmd.CommandPath() + "\" Requires exactly 1 arg.")
	}
	return cobra.OnlyValidArgs(cmd, args)
}
