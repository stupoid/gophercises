package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets DB by deleting the boltdb file",
	Run: func(cmd *cobra.Command, args []string) {
		ResetDB("/task/task.db")
		fmt.Println("DB Reset.")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
