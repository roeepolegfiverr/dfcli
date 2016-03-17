package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// queueCmd represents the queue command
var queueCmd = &cobra.Command{
	Use:   "queue <server name>",
	Short: "Join the waiting queue for a server",
	Long: `Join a Plike environment server waiting queue. You'll be notified when the server is free	to use.`,
}

func init() {
	RootCmd.AddCommand(queueCmd)
	queueCmd.PersistentPreRunE = persistentPreRunE
	queueCmd.PreRunE = queuePreRunE
	queueCmd.Run = queue
}

func queuePreRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Wrong number of arguments (expected 1).")
	}

	postData.ServerName = args[0]
	return nil
}

func queue(cmd *cobra.Command, args []string) {
	srv, err := Post("queue", postData.User.AuthToken, postData.toBytesArray())
	if err != nil {
		log.Fatal("%s%s%s\n", RED, err.Error(), NORMAL)
		return
	}

	fmt.Printf("Queued for %s\n", srv.Name)
}
