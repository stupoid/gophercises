/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a task on your TODO list",
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
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You have deleted the \"%s\" task.\n", t.Text)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
