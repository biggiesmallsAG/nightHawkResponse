package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	nhr "nighthawkresponse"
	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/auth"
	"nighthawkresponse/api/handlers/config"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	elastic "gopkg.in/olivere/elastic.v5"
)

const DEFAULT_PASSWORD = "nighthawk"
const NHINDEX = "nighthawk"
const NHACCOUNT = "accounts"

var (
	conf   *config.ConfigVars
	err    error
	client *elastic.Client
	query  elastic.Query
)

func init() {
	conf, err = config.ReadConfFile()
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to initialize config read")
		return
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.LogDebug(api.DEBUG, "Failed to initialize elastic client")
		return
	}
}

type UAMConfig struct {
	createAccount bool
	resetPassword bool
	Username      string
	Password      string
	Role          string
	Version       bool
}

func main() {
	var config UAMConfig

	flag.StringVar(&config.Username, "username", "", "Account username")
	flag.StringVar(&config.Password, "password", "", "Account password")
	flag.StringVar(&config.Role, "role", "User", "Account role")
	flag.BoolVar(&config.createAccount, "create-account", false, "Create new account only if account does not exist")
	flag.BoolVar(&config.resetPassword, "reset-password", false, "Reset account password to default")
	flag.BoolVar(&config.Version, "version", false, "Show version information")

	flag.Parse()

	if config.Version {
		nhr.ShowVersion("nhr-uam")
		os.Exit(0)
	}

	if config.Username == "" {
		fmt.Println("Account username is required")
		os.Exit(1)
	}

	if config.resetPassword {
		err := resetUserPassword(config.Username)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Password reset for username %s completed", config.Username)
		return
	}

	if config.createAccount {
		if config.Password == "" {
			err = createUserAccountInteractive(config.Username)
			if err != nil {
				fmt.Printf("Error creating user account - %s\n", err.Error())
			}
			return
		}

		err = createUserAccount(config.Username, config.Password, config.Role)
		if err != nil {
			fmt.Printf("Error creating user account - %s\n", err.Error())
		}

	}
}

func resetUserPassword(username string) error {
	validUser, userAcc := auth.UserExists(username)
	if !validUser {
		return errors.New(fmt.Sprintf("User %s does not exists\n", username))
	}
	passwordHash, _ := auth.HashPassword(DEFAULT_PASSWORD)
	out, err := client.Update().
		Index(NHINDEX).
		Type(NHACCOUNT).
		Id(userAcc.DocId).
		Doc(map[string]interface{}{"password_hash": passwordHash}).
		Do(context.Background())

	if err != nil {
		return err
	}

	if out.Result != "updated" {
		return errors.New(fmt.Sprintf("Password set activity completed with status %s\n", out.Result))
	}

	return nil
}

func createUserAccount(username, password, role string) error {
	validUser, userAcc := auth.UserExists(username)
	if validUser {
		return errors.New(fmt.Sprintf("Cannot create new user. User %s exists\n", username))
	}

	// Set default fields
	// and update as required
	userAcc.Default()
	userAcc.Username = strings.ToLower(username)
	userAcc.Role = strings.ToLower(role)
	userAcc.Password = password
	err = userAcc.Validation()
	if err != nil {
		return err
	}

	// Unset cleartext password value
	userAcc.Password = ""

	userAcc.PasswordHash, _ = auth.HashPassword(password)
	jsonAcc, _ := json.Marshal(userAcc)

	res, err := client.Index().
		Index(NHINDEX).
		Type(NHACCOUNT).
		BodyJson(string(jsonAcc)).
		Do(context.Background())

	if err != nil {
		return err
	}

	if res.Result != "created" {
		return errors.New(fmt.Sprintf("Error creating new user %s\n", username))
	}

	return nil
}

func createUserAccountInteractive(username string) error {
	var password, password2 string
	var role string

	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password = string(bytePassword)

	fmt.Print("\nConfirm Password: ")
	bytePassword, _ = terminal.ReadPassword(int(syscall.Stdin))
	password2 = string(bytePassword)

	if password != password2 {
		fmt.Println("Password did not match")
		os.Exit(2)
	}

	fmt.Printf("\nRole [admin/user]: ")
	fmt.Scanln(&role)

	return createUserAccount(username, password, role)
}
