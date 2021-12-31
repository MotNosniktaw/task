/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/motnosniktaw/task/database"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Sets a task as complete.",
	Long:  "Updates a task as completed. Provide list index to identify outstanding task.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Id of completed task must be provided.")
			os.Exit(1)
		}

		i, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task id format.")
			os.Exit(1)
		}

		database.MarkTaskAsCompleted(i)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
