package cmd

import (
	"errors"
	"fmt"
	"gophercises/task/db"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type fakeTask struct {
	err   error
	tasks []db.Task
}

func (f *fakeTask) createTask(s string) (int, error) {
	return 0, f.err
}

func TestAdd(t *testing.T) {
	var myCmd *cobra.Command
	home, _ := homedir.Dir()
	fmt.Println(home)
	dbPath := filepath.Join(home, "test.db")
	db.Init(dbPath)

	t.Run("it adds task successfully", func(t *testing.T) {
		addCmd.Run(myCmd, []string{"test_key"})
	})

	t.Run("it fails to create task if error occurs", func(t *testing.T) {
		f := &fakeTask{err: errors.New("Failed")}
		db.NewCreateTask = f.createTask
		addCmd.Run(myCmd, []string{"test_key"})
	})

	defer func() {
		db.NewCreateTask = db.CreateTask
	}()
}
