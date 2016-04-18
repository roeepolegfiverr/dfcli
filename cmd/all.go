// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"dfcli/auth"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Show all Plike servers",
	Long:  `Show all Plike servers`,
}

func init() {
	RootCmd.AddCommand(allCmd)

	allCmd.PersistentPreRunE = persistentPreRunE
	allCmd.Run = all
}

func all(cmd *cobra.Command, args []string) {
	envRes, err := Get("all", postData.User.AuthToken)
	if err != nil {
		log.Fatal("%s%s%s\n", RED, err.Error(), NORMAL)
		return
	}

	fmt.Printf("-- %s -------\n", envRes.Name)
	fmt.Printf("%-25s %-25s %-25s\n", "Server", "User", "Released At")
	for _, srv := range envRes.Servers {
		clr := NORMAL
		if srv.User == (auth.User{}) {
			clr = GREEN
		}
		fmt.Printf("%s%-25s %-25v %-25v%s\n", clr, srv.Name, srv.User.Name, srv.ReleaseDate, NORMAL)
	}
}
