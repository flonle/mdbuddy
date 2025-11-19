package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of MDBuddy",
	Long:  "All software has versions. This is MDBuddy's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MDBuddy v0.0.0 -- HEAD")
		// TODO add build info, get version from git tag, etc
	},
}
