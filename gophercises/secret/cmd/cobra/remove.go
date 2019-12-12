package cobra

import (
	"fmt"
	"gophercises/secret"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		key := args[0]
		err := v.Remove(key)
		if err != nil {
			fmt.Println("value not found")
			return
		}
		fmt.Println("Value removed successfully!")
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
