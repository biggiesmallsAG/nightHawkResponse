package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"
	"strings"

	"golang.org/x/crypto/bcrypt"
	elastic "gopkg.in/olivere/elastic.v5"
)

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

// bcrypt code below is taken from
// https://gowebexamples.com/password-hashing/

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Account struct {
	DocId        string `json:"doc_id,omitempty"`
	Firstname    string `json:"first_name,omitempty"`
	Lastname     string `json:"last_name,omitempty"`
	Username     string `json:"username"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email,omitempty"`
	Role         string `json:"role,omitempty"`
}

func (config *Account) Default() {
	config.Firstname = ""
	config.Lastname = ""
	config.Username = ""
	config.Password = ""
	config.PasswordHash = ""
	config.Email = ""
	config.Role = "User"
}

func (config *Account) LoadParams(data []byte) {
	config.Default()

	var tconfig Account
	json.Unmarshal(data, &tconfig)

	if tconfig.Firstname != "" {
		config.Firstname = tconfig.Firstname
	}
	if tconfig.Lastname != "" {
		config.Lastname = tconfig.Lastname
	}
	if tconfig.Username != "" {
		config.Username = tconfig.Username
	}
	if tconfig.Password != "" {
		config.Password = tconfig.Password
	}

	if tconfig.Email != "" {
		config.Email = tconfig.Email
	}
	if tconfig.Role != "" {
		config.Role = tconfig.Role
	}
}

func UserExists(username string) (bool, Account) {
	var acc Account

	query = elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("username.keyword", username))
	sr, err := client.Search().Index(NHINDEX).Type(NHACCOUNT).Query(query).Do(context.Background())
	if err != nil {
		api.LogError(api.DEBUG, err)
		return false, acc // If cannot determine assue it already exists
	}

	if sr.Hits.TotalHits >= 1 {

		for _, hit := range sr.Hits.Hits {
			json.Unmarshal(*hit.Source, &acc)
			acc.DocId = hit.Id
			break
		}
		return true, acc
	}

	return false, acc
}

func HttpErrorReturn(w http.ResponseWriter, r *http.Request, message string, err error) {
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] %s %s, %s", r.Method, r.RequestURI, err.Error()))
	fmt.Fprintln(w, api.HttpFailureMessage(message))
}

func HttpSuccessReturn(w http.ResponseWriter, r *http.Request, message string, hits int64) {
	api.LogDebug(api.DEBUG, fmt.Sprintf("[+] %s %s, %s", r.Method, r.RequestURI, message))
	fmt.Fprintln(w, api.HttpSuccessMessage("200", message, hits))
}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read HTTP Reuqest", err)
		return
	}

	var acc Account
	acc.LoadParams(body)

	conf, err = config.ReadConfFile()
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read config file", err)
		return
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		HttpErrorReturn(w, r, "Cannot connect to elasticsearch", err)
		return
	}

	boolUserExits, _ := UserExists(acc.Username)
	if boolUserExits {
		HttpErrorReturn(w, r, "User already exists", errors.New("User already exists"))
		return
	}

	// Compute BCrypt Hash
	acc.PasswordHash, err = HashPassword(acc.Password)
	if err != nil {
		fmt.Printf("Error hashingpassword, %s\n", err.Error())
	}

	// Set Password field empty so it is not indexed
	acc.Password = ""
	jsonAccount, _ := json.Marshal(acc)

	client.Index().Index(NHINDEX).Type(NHACCOUNT).BodyJson(string(jsonAccount)).Do(context.Background())
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read HTTP Reuqest", err)
		return
	}

	//fmt.Println("RequestHeaderToken ", r.Header.Get("NHR-TOKEN"))

	var passwdData map[string]string
	json.Unmarshal(body, &passwdData)
	passwdData["clientip"] = "127.0.0.1"
	passwdData["token"] = r.Header.Get("NHR-TOKEN")

	conf, err = config.ReadConfFile()
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read config file", err)
		return
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		HttpErrorReturn(w, r, "Cannot connect to elasticsearch", err)
		return
	}

	validSession, err := IsSessionTokenValid(passwdData["username"], passwdData["clientip"], passwdData["token"])
	if !validSession {
		HttpErrorReturn(w, r, "Session is not valid", err)
		return
	}

	// Get account information from Elasticsearch
	// Verify current password match and
	// update hashed password
	_, acc := UserExists(passwdData["username"])
	if CheckPasswordHash(passwdData["password"], acc.PasswordHash) {
		newPasswordHash, _ := HashPassword(passwdData["new_password"])
		out, err := client.Update().Index(NHINDEX).Type(NHACCOUNT).Id(acc.DocId).Doc(map[string]interface{}{"password_hash": newPasswordHash}).Do(context.Background())
		if err != nil {
			HttpErrorReturn(w, r, "Failed to change password", err)
			return
		}

		if out.Result != "updated" {
			HttpErrorReturn(w, r, "Failed to change password", errors.New(fmt.Sprintf("Failed to update password with result %s", out.Result)))
			return
		}
	}

	// default response is password changed successfully
	HttpSuccessReturn(w, r, "Password Changed", 1)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read HTTP Reuqest", err)
		return
	}

	conf, err = config.ReadConfFile()
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read config file", err)
		return
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		HttpErrorReturn(w, r, "Cannot connect to elasticsearch", err)
		return
	}

	var postData map[string]string
	json.Unmarshal(body, &postData)

	validUser, acc := UserExists(postData["delete_account"])
	if !validUser {
		HttpErrorReturn(w, r, "User not found", errors.New("User not found"))
		return
	}
	res, err := client.Delete().Index(NHINDEX).Type(NHACCOUNT).Id(acc.DocId).Do(context.Background())

	if err != nil || res.Result != "deleted" {
		HttpErrorReturn(w, r, "Error deleting user", err)
		return
	}

	HttpSuccessReturn(w, r, "Account deleted", 1)
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	/// verify requesting user has admin role
	token := r.Header.Get("NHR-TOKEN")
	adminUser := GetUseranmeByToken(token)
	if adminUser == "" {
		HttpErrorReturn(w, r, "Could not find valid session for admin user", errors.New("Admin session is not valid"))
		return
	}

	adminExists, adminAcc := UserExists(adminUser)
	if !adminExists {
		HttpErrorReturn(w, r, fmt.Sprintf("User account %s does not exits", adminUser), errors.New(fmt.Sprintf("User account %s does not exits", adminUser)))
		return
	}

	if strings.ToLower(adminAcc.Role) != "admin" {
		HttpErrorReturn(w, r, fmt.Sprintf("User account %s does not admin user", adminUser), errors.New(fmt.Sprintf("User account %s is not admin user", adminUser)))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read HTTP Reuqest", err)
		return
	}

	var postData map[string]string
	json.Unmarshal(body, &postData)

	validUser, acc := UserExists(postData["username"])
	if !validUser {
		HttpErrorReturn(w, r, "User not found", errors.New("User not found"))
		return
	}

	newPasswordHash, _ := HashPassword(postData["new_password"])
	out, err := client.Update().Index(NHINDEX).Type(NHACCOUNT).Id(acc.DocId).Doc(map[string]interface{}{"password_hash": newPasswordHash}).Do(context.Background())
	if err != nil || out.Result != "updated" {
		HttpErrorReturn(w, r, "Failed to change password", err)
		return
	}
	HttpSuccessReturn(w, r, "Password Set", 1)
}
