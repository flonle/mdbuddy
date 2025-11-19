package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/flonle/mdbuddy/renderer"
	"github.com/spf13/cobra"
)

func init() {
	renderCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
	rootCmd.AddCommand(renderCmd)
}

var renderCmd = &cobra.Command{
	Use:   "render [file]",
	Short: "Render a single markdown file",
	Long: `Render a single markdown file to HTML and print to stdout.

This is useful for one-off rendering or integration with other tools.

Examples:
	mdbuddy render README.md
	mdbuddy render docs/guide.md > output.html
	echo "# Hello" | mdbuddy render
	mdbuddy render input.md --output result.html
`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRender,
}

func runRender(cmd *cobra.Command, args []string) error {
	var input []byte
	var inputFile string
	var err error
	if len(args) == 0 {
		input, err = readStdin()
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else {
		inputFile = args[0]
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", inputFile)
		}
		input, err = os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", inputFile, err)
		}
	}

	// Get output file
	var outputFile io.Writer
	outputFileName, _ := cmd.Flags().GetString("output")
	if outputFileName != "" {
		outputFile, err = os.Create(outputFileName)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	} else {
		outputFile = os.Stdout
	}

	// Render
	renderer.Render(input, outputFile)

	// Some extra info on stdin, if it isn't used to print the HTML
	if outputFile != os.Stdout {
		if inputFile != "" {
			fmt.Printf("✅ Rendered %s → %s\n", inputFile, outputFileName)
		} else {
			fmt.Printf("✅ Rendered stdin → %s\n", outputFileName)
		}
	}

	return nil
}

func readStdin() ([]byte, error) {
	// Check if stdin has data
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	// Check if stdin is from a pipe/redirect
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return io.ReadAll(os.Stdin)
	}

	// Interactive mode - prompt user
	fmt.Fprint(os.Stderr, "Enter markdown content (Ctrl+D to finish):\n")
	return io.ReadAll(os.Stdin)
}
