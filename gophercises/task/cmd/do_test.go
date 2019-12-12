package cmd

import (
	"errors"
	"gophercises/task/db"
	"testing"

	"github.com/spf13/cobra"
)

func (f *fakeTask) deleteTask(k int) error {
	return f.err
}

func (f *fakeTask) allTask() ([]db.Task, error) {
	return f.tasks, f.err
}

func TestDo(t *testing.T) {
	var myCmd *cobra.Command

	t.Run("it adds task successfully", func(t *testing.T) {
		addCmd.Run(myCmd, []string{"test_key"})
		doCmd.Run(myCmd, []string{"1"})
	})

	t.Run("it fails to create task if error occurs", func(t *testing.T) {
		f := &fakeTask{err: errors.New("Failed")}
		db.NewDeleteTask = f.deleteTask
		doCmd.Run(myCmd, []string{"1"})
	})

	t.Run("it fails if task is not present", func(t *testing.T) {
		doCmd.Run(myCmd, []string{"100"})
	})

	t.Run("it fails if all task is having error", func(t *testing.T) {
		f := &fakeTask{err: errors.New("Failed")}
		db.NewAllTasks = f.allTask
		doCmd.Run(myCmd, []string{"1"})
	})

	t.Run("it fails if task is not provided", func(t *testing.T) {
		doCmd.Run(myCmd, []string{""})
	})

	defer func() {
		db.NewCreateTask = db.CreateTask
		db.NewAllTasks = db.AllTasks
	}()
}
