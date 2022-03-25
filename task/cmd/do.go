package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		path, err := getDbPath()
		if err != nil {
			log.Fatal(err)
		}
		t, err := PopTask(path, BucketName, incompleteTasksKey, id-1)
		if err != nil {
			log.Fatal(err)
		}
		t.UpdatedAt = time.Now()
		err = putTask(path, BucketName, completedTasksKey, t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You have completed the \"%s\" task.\n", t.Text)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
