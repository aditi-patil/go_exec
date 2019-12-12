package cobra

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRemove(t *testing.T) {
	var myCmd *cobra.Command

	t.Run("it gives vakue of key if key is present", func(t *testing.T) {
		setCmd.Run(myCmd, []string{"twit_api1", "newvalue"})
		removeCmd.Run(myCmd, []string{"twit_api1"})
	})

	t.Run("it gives vakue of key if key is present", func(t *testing.T) {
		removeCmd.Run(myCmd, []string{"twit_api12"})
	})
}
