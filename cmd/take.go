package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var takeCmd = &cobra.Command{
	Use:   "take <server_name>",
	Short: "Take a server",
	Long:  "Take a server from the Plike environment for one hour (can be extended).",
}

func init() {
	RootCmd.AddCommand(takeCmd)

	takeCmd.PersistentPreRunE = persistentPreRunE
	takeCmd.PreRunE = takePreRunE
	takeCmd.Run = take
}

func takePreRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Wrong number of arguments (expected 1).")
	}

	// TODO: Validate - can't take more than one server
	postData.ServerName = args[0]
	postData.ReleaseDate = time.Now().Add(time.Hour).Format(time.RFC3339)
	return nil
}

func take(cmd *cobra.Command, args []string) {
	srv, err := Post("take", postData.User.AuthToken, postData.toBytesArray())
	if err != nil {
		log.Fatal("%s%s%s\n", RED, err.Error(), NORMAL)
		return
	}

	fmt.Printf("%s - %s was taken until %s\n", srv.Environment, srv.Name, srv.ReleaseDate)
}
