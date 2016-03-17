package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release <server name>",
	Short: "Release a Plike server.",
	Long:  `Release a Plike server. Once released others may overwrite anything left on that machine.`,
}

func init() {
	RootCmd.AddCommand(releaseCmd)
	releaseCmd.PersistentPreRunE = persistentPreRunE
	releaseCmd.PreRunE = releasePreRunE
	releaseCmd.Run = release
}

func releasePreRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Wrong number of arguments (expected 1).")
	}

	postData.ServerName = args[0]
	return nil
}

func release(cmd *cobra.Command, args []string) {

	if _, err := Post("release", postData.User.AuthToken, postData.toBytesArray()); err != nil {
		log.Fatal("%s%s%s\n", RED, err.Error(), NORMAL)
		return
	}

	fmt.Printf("Released %s!\n", postData.ServerName)
}
