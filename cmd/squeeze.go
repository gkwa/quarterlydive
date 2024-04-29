package cmd

import (
	"github.com/spf13/cobra"
	"github.com/taylormonacelli/quarterlydive/squeeze"
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
		squeeze.RunSqueezer(directory)
	},
}

func init() {
	rootCmd.AddCommand(squeezeCmd)
}
