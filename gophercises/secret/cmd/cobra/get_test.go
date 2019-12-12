package cobra

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestGet(t *testing.T) {
	var myCmd *cobra.Command
	t.Run("it gives value of key if key is present", func(t *testing.T) {
		setCmd.Run(myCmd, []string{"twit_api1", "newvalue"})
		getCmd.Run(myCmd, []string{"twit_api1"})
	})

	t.Run("it gives value of key if key is present", func(t *testing.T) {
		getCmd.Run(myCmd, []string{"twit_api2"})
	})
}
