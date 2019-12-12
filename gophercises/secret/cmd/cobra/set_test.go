package cobra

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSet(t *testing.T) {

	t.Run("set key and value", func(t *testing.T) {
		var myCmd *cobra.Command
		args := []string{"twit_api", "newvalue"}
		setCmd.Run(myCmd, args)
	})

}
