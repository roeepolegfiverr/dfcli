package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"time"
)

var extendCmd = &cobra.Command{
	Use:   "extend <server name> <add X hours>",
	Short: "Extend server usage by X hours.",
	Long:  "Extend taken server from the Plike environment by X hours.",
}

func init() {
	RootCmd.AddCommand(extendCmd)
	extendCmd.PersistentPreRunE = persistentPreRunE
	extendCmd.PreRunE = extendPreRunE
	extendCmd.Run = extend
}

func extendPreRunE(cmd *cobra.Command, args []string) (err error) {
	if len(args) != 2 {
		return errors.New("Wrong number of arguments (expected 2).")
	}

	postData.ServerName = args[0]
	hours, _ := strconv.Atoi(args[1])
	postData.ReleaseDate = time.Now().Add(time.Hour * time.Duration(hours)).Format(time.RFC3339)

	return nil
}

func extend(cmd *cobra.Command, args []string) {

	srv, err := Post("extend", postData.User.AuthToken, postData.toBytesArray())
	if err != nil {
		log.Fatal("%s%s%s\n", RED, err.Error(), NORMAL)
		return
	}
	fmt.Println(srv)
	fmt.Printf("%s - %s was extended until %s\n", srv.Environment, srv.Name, srv.ReleaseDate)
}
