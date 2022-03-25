package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// completedCmd represents the completed command
var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all of your completed tasks for today",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := getTasks(completedTasksKey)
		if err != nil {
			log.Fatal(err)
		}
		var tasksCompletedToday []Task
		for _, t := range tasks {
			if SameDate(t.UpdatedAt, time.Now()) {
				tasksCompletedToday = append(tasksCompletedToday, t)
			}
		}

		if len(tasksCompletedToday) > 0 {
			fmt.Println("You have finished the following tasks today:")
			for _, t := range tasksCompletedToday {
				fmt.Printf("- %s\n", t.Text)
			}
		} else {
			fmt.Println("You have no completed tasks.")
		}
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
