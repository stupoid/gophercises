package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("No task id specified")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		t, err := DoTask(id - 1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You have completed the \"%s\" task.\n", t.Text)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
