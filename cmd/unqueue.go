package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// unqueueCmd represents the unqueue command
var unqueueCmd = &cobra.Command{
	Use:   "unqueue",
	Short: "Remove your name form a Plike server's waiting queue.",
	Long:  `Remove your name from a Plike server's waiting queue. Once removed no notification will be sent when the server is released.`,
}

func init() {
	RootCmd.AddCommand(unqueueCmd)

	unqueueCmd.PersistentPreRunE = persistentPreRunE
	unqueueCmd.PreRunE = unqueuePreRunE
	unqueueCmd.Run = unqueue
}

func unqueuePreRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Wrong number of arguments (expected 1).")
	}

	postData.ServerName = args[0]
	return nil
}

func unqueue(cmd *cobra.Command, args []string) {

	if _, err := Post("unqueue", postData.User.AuthToken, postData.toBytesArray()); err != nil {
		log.Fatal("%s%s%s\n", RED, err.Error(), NORMAL)
		return
	}

	fmt.Printf("Unqueued for %s!\n", postData.ServerName)
}
