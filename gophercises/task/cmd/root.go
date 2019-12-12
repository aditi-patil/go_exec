package cmd

import "github.com/spf13/cobra"

// RootCmd root command of CLI manager
var RootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a CLI task manager",
}
