package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

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

// UserExxists verifies user exists in database
// and return user Account data
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

// ChangePassword changed currently logged on user password
// api_uri: POST /api/v1/user/password/change
// post_data: {"password":"currentpassword", "new_password":"newpassword123"}
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	isauth, message := IsAuthenticatedSession(w, r)
	if !isauth {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	var passwdData map[string]string
	json.Unmarshal(body, &passwdData)

	// Get username from authenticated sesion table
	passwdData["token"] = r.Header.Get("NHR-TOKEN")
	passwdData["username"] = GetUseranmeByToken(passwdData["token"])

	// Get account information from Elasticsearch
	// Verify current password match and update hashed password
	_, acc := UserExists(passwdData["username"])
	if CheckPasswordHash(passwdData["password"], acc.PasswordHash) {
		newPasswordHash, _ := HashPassword(passwdData["new_password"])
		out, err := client.Update().Index(NHINDEX).Type(NHACCOUNT).Id(acc.DocId).Doc(map[string]interface{}{"password_hash": newPasswordHash}).Do(context.Background())
		if err != nil || out.Result != "updated" {
			api.HttpResponseReturn(w, r, "failed", "Failed to changed password", nil)
			return
		}
	}

	// default response is password changed successfully
	api.HttpResponseReturn(w, r, "success", "Password change completed", passwdData["username"])
}

// CreateNewUser creates a new user
// This operation is only allowed to authenticated admin user
// api_uri: POST /api/v1/admin/user/create
// post_data: {"username":"user", "password":"password123", "role":"user"}
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// check if requesting user is admin user
	isadmin, message := IsAuthenticatedAdminSession(w, r)
	if !isadmin {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	var acc Account
	acc.LoadParams(body)

	boolUserExits, _ := UserExists(acc.Username)
	if boolUserExits {
		api.HttpResponseReturn(w, r, "failed", "User already exits", acc.Username)
		return
	}

	// Compute BCrypt Hash
	acc.PasswordHash, err = HashPassword(acc.Password)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", "Internal error hashing password", nil)
	}

	// Set Password field empty so it is not indexed
	acc.Password = ""
	jsonAccount, _ := json.Marshal(acc)

	res, err := client.Index().Index(NHINDEX).Type(NHACCOUNT).BodyJson(string(jsonAccount)).Do(context.Background())

	if err != nil || !res.Created {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	resdata := make(map[string]string)
	resdata["username"] = acc.Username
	resdata["role"] = acc.Role
	api.HttpResponseReturn(w, r, "success", "User account created", "")
}

// DeleteUser deletes useraccount
// This operation is only allowed to authenticated admin user
// api_uri: POST /api/v1/admin/user/delete
// post_data: {"username":"user"}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// check if requesting user is admin user
	isadmin, message := IsAuthenticatedAdminSession(w, r)
	if !isadmin {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	var postData map[string]string
	json.Unmarshal(body, &postData)

	validUser, acc := UserExists(postData["username"])
	if !validUser {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}
	res, err := client.Delete().Index(NHINDEX).Type(NHACCOUNT).Id(acc.DocId).Do(context.Background())

	if err != nil || res.Result != "deleted" {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}
	api.HttpResponseReturn(w, r, "success", "User account delete", postData)
}

// SetPassword sets password for a user account.
// This operation is only allowed to authenticated administrator user
// api_uri: POST /api/v1/admin/password/set
// post_data: {"username", "user", "new_password": "password123"}
func SetPassword(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// check if requesting user is admin user
	isadmin, message := IsAuthenticatedAdminSession(w, r)
	if !isadmin {
		api.HttpResponseReturn(w, r, "failed", message, nil)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	var postData map[string]string
	json.Unmarshal(body, &postData)

	validUser, acc := UserExists(postData["username"])
	if !validUser {
		api.HttpResponseReturn(w, r, "failed", "User not found", postData["username"])
		return
	}

	newPasswordHash, _ := HashPassword(postData["new_password"])
	out, err := client.Update().Index(NHINDEX).Type(NHACCOUNT).Id(acc.DocId).Doc(map[string]interface{}{"password_hash": newPasswordHash}).Do(context.Background())
	if err != nil || out.Result != "updated" {
		api.HttpResponseReturn(w, r, "failed", err.Error(), postData["username"])
		return
	}
	api.HttpResponseReturn(w, r, "success", "Password set completed", postData["username"])
}
