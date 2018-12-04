package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints mfaws version information",
	Long:  `Prints mfaws version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v, commit %v, built at %v\n", VERSION, COMMIT, DATE)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
