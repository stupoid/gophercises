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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("No task id specified")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		// cli displays list starting from 1
		t, err := PopTask(incompleteTasksKey, id-1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("You have deleted the \"%s\" task.\n", t.Text)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
