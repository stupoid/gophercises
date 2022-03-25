package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets DB by deleting the boltdb file",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getDbPath()
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Removed DB at \"%s\".\n", path)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
