package cmd

import (
	"bytes"
	"dfcli/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	RED    = "\x1b[31;1m"
	NORMAL = "\x1b[0m"
	GREEN  = "\x1b[32;1m"
	PLIKE  = "Plike"
)

// Data Types
type EnvResponse struct {
	Name    string   `json:"name"`
	Servers []Server `json:"servers"`
}

type Server struct {
	Environment string    `json:"environment"`
	Name        string    `json:"name"`
	User        auth.User `json:"user"`
	ReleaseDate string    `json:"release_date"`
}

func (s *Server) isFree() bool {
	return s.User == (auth.User{})
}

type PostData struct {
	Environment string    `json:"environment"`
	ServerName  string    `json:"name"`
	User        auth.User `json:"user"`
	ReleaseDate string    `json:"release_date"`
}

func (p *PostData) toBytesArray() []byte {
	j, _ := json.Marshal(p)
	return []byte(j)
}

var postData PostData

func Post(cmdName, authToken string, payloadStr []byte) (srv Server, err error) {

	url := "http://" + os.Getenv("DFCLI_END_POINT") + "/server/rest/" + cmdName
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authToken))

	client := &http.Client{}
	rawResponse, err := client.Do(req)
	if err != nil {
		return srv, err
	}
	defer rawResponse.Body.Close()

	if rawResponse.StatusCode != 200 && rawResponse.StatusCode != 201 {
		return srv, errors.New(rawResponse.Status)
	}

	body, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return srv, err
	}

	err = json.Unmarshal(body, &srv)
	if err != nil {
		return srv, err
	}
	// fmt.Println(srv)
	return srv, nil
}

func Get(cmdName, authToken string) (res EnvResponse, err error) {
	url := "http://" + os.Getenv("DFCLI_END_POINT") + "/server/rest/" + cmdName + "?" + fmt.Sprintf("environment=%s", PLIKE)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authToken))

	client := &http.Client{}
	rawResponse, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer rawResponse.Body.Close()

	if rawResponse.StatusCode != 200 {
		return res, errors.New(rawResponse.Status)
	}

	body, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func persistentPreRunE(cmd *cobra.Command, args []string) (err error) {
	user, err := auth.ReadAuth()
	if err != nil {
		fmt.Println("Error reading authentication data. Please run:\n dfcli init")
		return err
	}

	postData = PostData{Environment: PLIKE, User: user}
	return nil
}
