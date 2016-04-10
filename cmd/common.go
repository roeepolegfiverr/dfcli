package cmd

import (
	"bytes"
	"crypto/tls"
	"dfcli/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	RED                 = "\x1b[31;1m"
	NORMAL              = "\x1b[0m"
	GREEN               = "\x1b[32;1m"
	PLIKE               = "Plike"
	REQUEST_TIMEOUT_SEC = 5
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
var retriesCount = 0
var retry = true

func Post(cmdName, authToken string, payloadStr []byte) (srv Server, err error) {

	url := "https://" + os.Getenv("DFCLI_END_POINT") + "/server/rest/" + cmdName

	body, err := doRequest("POST", url, authToken, payloadStr)
	if err != nil {
		return srv, err
	}

	err = json.Unmarshal(body, &srv)
	if err != nil {
		return srv, err
	}

	return srv, nil
}

func Get(cmdName, authToken string) (res EnvResponse, err error) {
	url := "https://" + os.Getenv("DFCLI_END_POINT") + "/server/rest/" + cmdName + "?" + fmt.Sprintf("environment=%s", PLIKE)

	body, err := doRequest("GET", url, authToken, nil)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func doRequest(httpMethod, url, authToken string, payloadStr []byte) (res []byte, err error) {
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(payloadStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+authToken)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(REQUEST_TIMEOUT_SEC * time.Second),
	}

	rawResponse := &http.Response{}
	for retriesCount < MaxRetries && retry {

		rawResponse, err = client.Do(req)

		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("Timeout, retrying...", retriesCount)
			retriesCount++

		} else if err != nil {
			return res, err

		} else if err == nil {
			defer rawResponse.Body.Close()
			retry = false
		}
	}

	//
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return res, errors.New("Timeout!")
	}

	if rawResponse.StatusCode != 200 {
		return res, errors.New(rawResponse.Status)
	}

	res, err = ioutil.ReadAll(rawResponse.Body)
	return res, err
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
