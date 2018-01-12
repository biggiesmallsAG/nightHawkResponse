package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	api "nighthawkapi/api/core"
	"nighthawkapi/api/handlers/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Login function authenticates user
// api_uri: /api/v1/auth/user/login
// post_data: {"username":"user1", "password":"password123"}
func Login(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", "Failed to read HTTP request", nil)
		return
	}

	loginData := make(map[string]string)
	json.Unmarshal(body, &loginData)

	conf, err = config.ReadConfFile()
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", conf.ServerHttpScheme(), conf.ServerHost(), conf.ServerPort())))
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	// Account informatino extracted from Server
	boolUserExits, sacc := UserExists(loginData["username"])
	if !boolUserExits {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	passwordMatched := CheckPasswordHash(loginData["password"], sacc.PasswordHash)
	if passwordMatched {
		clientip := "127.0.0.1" // TODO: Extract from request header
		sessionToken, err := CreateSession(loginData["username"], sacc.Role, clientip)
		if err != nil {
			api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
			return
		}
		w.Header().Set("NHR-TOKEN", sessionToken)

		loginData["token"] = sessionToken
		delete(loginData, "password")
		api.HttpResponseReturn(w, r, "success", "Authentication completed", loginData)
	} else {
		api.HttpResponseReturn(w, r, "failed", "Password did not match", nil)
	}
}

// Logout function logouts authenticated user
// api_uri: POST /api/v1/auth/user/logout
// post_data: {"username": "user1"}
func Logout(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), nil)
		return
	}

	postData := make(map[string]string)
	json.Unmarshal(body, &postData)

	postData["token"] = r.Header.Get("NHR-TOKEN")
	postData["clientip"] = "127.0.0.1"

	_, err = DestroySession(postData["username"], postData["clientip"], postData["token"])
	if err != nil {
		api.HttpResponseReturn(w, r, "failed", err.Error(), postData)
		return
	}

	api.HttpResponseReturn(w, r, "success", "Logout completed", postData)
}
