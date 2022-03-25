package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getDbPath()
		if err != nil {
			log.Fatal(err)
		}
		tasks, err := getTasks(path, BucketName, incompleteTasksKey)
		if err != nil {
			log.Fatal(err)
		}
		if len(tasks) > 0 {
			fmt.Println("You have the following tasks:")
			for i, t := range tasks {
				fmt.Printf("%d. %s\n", i+1, t.Text)
			}
		} else {
			fmt.Println("You have no incomplete tasks.")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
