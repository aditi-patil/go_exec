package db

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

type fakeDB struct {
	err error
}

func (f *fakeDB) dbView() ([]Task, error) {
	return nil, f.err
}

func (f *fakeDB) dbUpdate(task string) (int, error) {
	return 0, f.err
}

func TestInit(t *testing.T) {

	t.Run("it opens db connection if error is not present", func(t *testing.T) {
		home, _ := homedir.Dir()
		dbPath := filepath.Join(home, "test.db")
		err := Init(dbPath)
		if err != nil {
			t.Errorf("failed to open db with error %v", err)
		}
	})

	t.Run("it returns error if db connection failed to open", func(t *testing.T) {
		home, _ := homedir.Dir()
		dbPath := filepath.Join(home, "testing/test.db")
		err := Init(dbPath)
		expected := "open " + dbPath + ": no such file or directory"
		assert.Equal(t, expected, err.Error())
	})

}

func TestCreateTask(t *testing.T) {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "test_create.db")
	Init(dbPath)
	t.Run("it creates a new task", func(t *testing.T) {
		id, _ := CreateTask("test_key")
		assert.NotNil(t, id)
	})

	t.Run("it returns error if boltdb failed to create task", func(t *testing.T) {
		f := &fakeDB{err: errors.New("Failed")}
		newDbUpdate = f.dbUpdate
		id, err := CreateTask("demo_key")
		assert.NotNil(t, err)
		assert.Equal(t, id, -1)
	})
}

func TestAllTasks(t *testing.T) {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "all_tasks.db")
	Init(dbPath)
	t.Run("it gives all tasks", func(t *testing.T) {
		tasks, _ := AllTasks()
		assert.NotEmpty(t, tasks)
	})

	t.Run("it returns error if boltdb failed to give tasks", func(t *testing.T) {
		f := &fakeDB{err: errors.New("Failed")}
		newDbView = f.dbView
		_, err := AllTasks()
		assert.NotNil(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "delete_tasks.db")
	Init(dbPath)
	t.Run("it gives all tasks", func(t *testing.T) {
		err := DeleteTask(1)
		assert.Nil(t, err)
	})

	t.Run("it returns error if boltdb failed to delete task", func(t *testing.T) {

	})
}
