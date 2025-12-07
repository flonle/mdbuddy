package cmd

import (
	"fmt"

	"github.com/flonle/mdbuddy/server"
	"github.com/spf13/cobra"
)

func init() {
	watchCmd.Flags().StringP("bind", "b", "", "Bind to this address (default: all interfaces)")
	watchCmd.Flags().StringP("port", "p", "", "Bind to this port (default: 3000)")
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch [files]...",
	Short: "Watch the given files or directories",
	Long: `Watch a set of files and/or directories, and continuously re-render a preview of the last changed markdown file in that set.
Supplying a directory is the same as supplying all files in it. Files without the .md extension are ignored. Linux only ¯\_(ツ)_/¯`,
	Example: `  mdbuddy watch README.md README2.md README3.md
  mdbuddy watch .`,
	Args: cobra.MinimumNArgs(1),
	RunE: runWatch,
}

func runWatch(cmd *cobra.Command, args []string) error {
	bind, _ := cmd.Flags().GetString("bind") // If not set, we get "", which is fine
	port, _ := cmd.Flags().GetString("port")
	if port == "" {
		port = "3000"
	}

	return server.ServePreview(fmt.Sprintf("%s:%s", bind, port), args)
}
