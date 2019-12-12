package cmd

import (
	"errors"
	"gophercises/task/db"
	"testing"

	"github.com/spf13/cobra"
)

func TestList(t *testing.T) {
	var myCmd *cobra.Command

	t.Run("it lists all tasks", func(t *testing.T) {
		listCmd.Run(myCmd, []string{})
	})

	t.Run("it fails if all task is having error", func(t *testing.T) {
		f := &fakeTask{err: errors.New("Failed")}
		db.NewAllTasks = f.allTask
		listCmd.Run(myCmd, nil)
	})

	t.Run("it returns if task is not present", func(t *testing.T) {
		f := &fakeTask{tasks: nil}
		db.NewAllTasks = f.allTask
		listCmd.Run(myCmd, nil)
	})
}
