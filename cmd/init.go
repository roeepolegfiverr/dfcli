package cmd

import (
	"bufio"
	"dfcli/auth"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup Devflow authentication data",
	Long:  `Setup Devflow authentication data: LDAP password and full user data (name, email and image url).`,
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Run = initRun
}

func initRun(cmd *cobra.Command, args []string) {
	var ldapName, ldapPass, email, imageUrl string

	fmt.Println("Enter the following details:")
	fmt.Printf("LDAP Name: ")
	fmt.Scanf("%s", &ldapName)

	fmt.Printf("LDAP Password: ")
	fmt.Scanf("%s", &ldapPass)

	fmt.Printf("DevFlow Image URL: ")
	fmt.Scanf("%s", &imageUrl)

	fmt.Printf("DevFlow Email: ")
	fmt.Scanf("%s", &email)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("DevFlow Username: ")
	userName, _ := reader.ReadString('\n')
	userName = strings.Replace(userName, "\n", "", -1)

	err := auth.SaveAuth(ldapName, ldapPass, userName, email, imageUrl)
	if err != nil {
		log.Fatal("%sERROR - SaveAuth - %s%s\n", RED, err.Error(), NORMAL)
		return
	}

	fmt.Println("Done!")
}
