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
The resulting HTML is completely self-contained.`,
	Example: `  mdbuddy render README.md
  mdbuddy render docs/guide.md > output.html
  echo "# Hello" | mdbuddy render
  mdbuddy render input.md --output result.html`,
	Args: cobra.ExactArgs(1),
	RunE: runRender,
}

func runRender(cmd *cobra.Command, args []string) error {
	var err error

	// Get input
	var input []byte
	var inputFile string
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

	// Get output
	var output io.Writer
	outputFileName, _ := cmd.Flags().GetString("output")
	if outputFileName != "" {
		output, err = os.Create(outputFileName)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	} else {
		output = os.Stdout
	}

	// Render input to output
	renderer.RenderBareNote(input, output)

	// Some extra info on stdin, if it isn't already used to print the HTML
	if output != os.Stdout {
		if inputFile != "" {
			fmt.Printf("✅ Rendered %s → %s\n", inputFile, outputFileName)
		} else {
			fmt.Printf("✅ Rendered stdin → %s\n", outputFileName)
		}
	}

	return nil
}

// Read & return stdin as []byte
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
