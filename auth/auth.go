package auth

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/lunixbochs/go-keychain"
	"os"
	"strings"
)

const (
	SERVICE_KEY   = "dfcli"
	USER_CSV_FILE = ".dfcli"
	RED           = "\x1b[31;1m"
	NORMAL        = "\x1b[0m"
	GREEN         = "\x1b[32;1m"
)

type User struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	AuthToken string
}

func (u *User) Blank() bool {
	return u.Email == "" && u.Image == "" && u.Email == ""
}

func SaveAuth(ldapName, ldapPass, username, email, imageUrl string) (err error) {
	csv := fmt.Sprintf("%s,%s,%s,%s", ldapName, username, email, imageUrl)

	file, err := os.Create(USER_CSV_FILE)
	if err != nil {
		fmt.Printf("%sERROR - SaveAuth - %s%s\n", RED, err.Error(), NORMAL)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(csv)
	if err != nil {
		fmt.Printf("%sERROR - SaveAuth - %s%s\n", RED, err.Error(), NORMAL)
		return err
	}
	file.Sync()

	err = saveToKeychain(ldapName, ldapPass)
	if err != nil {
		fmt.Printf("%sERROR - SaveAuth - %s%s\n", RED, err.Error(), NORMAL)
		return err
	}

	return nil
}

func ReadAuth() (user User, err error) {

	file, err := os.Open(USER_CSV_FILE)
	if err != nil {
		return user, err
	}
	defer file.Close()

	file_data := make([]byte, 1000)
	count, err := file.Read(file_data)
	if err != nil {
		return user, err
	}

	splited := strings.Split(string(file_data[:count]), ",")
	token, err := readKeychain(splited[0]) //ldapName
	if err != nil {
		return user, err
	}

	user = User{Name: splited[1], Email: splited[2], Image: splited[3], AuthToken: token}

	return user, err
}

func saveToKeychain(ldapName, ldapPass string) error {
	str := fmt.Sprintf("%s:%s", ldapName, ldapPass)
	b64String := b64.StdEncoding.EncodeToString([]byte(str))

	return keychain.Add(SERVICE_KEY, ldapName, b64String)
}

func readKeychain(ldapName string) (key string, err error) {
	return keychain.Find(SERVICE_KEY, ldapName)
}
