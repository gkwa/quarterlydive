package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylormonacelli/quarterlydive/squeeze"
)

var (
	assumeYes bool
	maxFiles  int
)

var squeezeCmd = &cobra.Command{
	Use:   "squeeze",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			err := cmd.Usage()
			if err != nil {
				panic(err)
			}
			return
		}
		directory := args[0]

		candidateFiles := squeeze.CountCandidateFiles(directory)
		if candidateFiles > maxFiles && !assumeYes {
			fmt.Printf("Found %d candidate files. Continue? [y/N]: ", candidateFiles)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if input != "y\n" && input != "Y\n" {
				fmt.Println("Aborted.")
				return
			}
		}

		squeeze.RunSqueezer(directory)
	},
}

func init() {
	rootCmd.AddCommand(squeezeCmd)

	squeezeCmd.Flags().BoolVar(&assumeYes, "assume-yes", false, "Automatically answer yes to prompts")
	squeezeCmd.Flags().IntVar(&maxFiles, "max-files", 20, "Maximum number of files to process without confirmation")
}
