/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/motnosniktaw/task/database"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all outstandings tasks.",
	Long:  "Display a list of all the outstanding tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		tasks := database.GetTasks()

		outstandingTasksCount := 0

		for _, t := range tasks {
			if !t.Complete {
				fmt.Printf("%d: %s\n", t.ID, t.Task)
				outstandingTasksCount++
			} else if all {
				fmt.Printf("%d: %s - COMPLETED\n", t.ID, t.Task)
			}
		}

		if outstandingTasksCount == 0 {
			fmt.Println("\nNo outstanding tasks.")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	listCmd.PersistentFlags().Bool("all", false, "Include completed tasks.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
