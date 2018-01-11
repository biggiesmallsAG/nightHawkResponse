package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"nighthawkapi/api/handlers/config"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

type AuthSession struct {
	Username        string
	ClientIp        string
	LoginTime       int64
	Token           string
	TokenHardExpiry int64 // after 24 hours
	IdleTimeOut     int64 // 30 minutes
	LastActivity    int64
}

func (config *AuthSession) NewSession() {
	config.LoginTime = time.Now().UTC().Unix()
	config.TokenHardExpiry = config.LoginTime + int64(86400) // expire after 24 hours
	config.IdleTimeOut = 1800

	//// Generating Token
	b := make([]byte, 32)
	rand.Read(b)
	config.Token = base64.URLEncoding.EncodeToString(b)
}

//=================================================
// AuthSession Handlers
//=================================================
var authsession []AuthSession // Store all session
var asmap map[string]int      // Store index to authsession

func init() {
	asmap = make(map[string]int)
}

func CreateSession(username string, clientip string) (string, error) {
	var as AuthSession
	as.Username = username
	as.ClientIp = clientip
	as.NewSession()

	authsession = append(authsession, as)
	asmap[as.Token] = len(authsession) - 1

	return as.Token, nil
}

func DestroySession(username string, clientip string, token string) (bool, error) {
	if len(authsession) == 0 {
		return false, errors.New("No active session found")
	}

	asIndex := asmap[token]

	//// Check all the failing conditions
	if asIndex < 0 || asIndex > len(authsession) {
		return false, errors.New("Cannot find valid index in asmap")
	}

	if authsession[asIndex].Username != username {
		return false, errors.New("Username does not match in AuthSession")
	}

	if authsession[asIndex].Token != token {
		return false, errors.New("Token does not match in AuthSession")
	}

	if authsession[asIndex].ClientIp != clientip {
		return false, errors.New("ClientIP does not match in AuthSession")
	}

	//// re-organize authsession and asmap
	tas := authsession
	authsession = authsession[:0] // clear slice
	asmap = nil                   // clear map
	asmap = make(map[string]int)
	j := 0
	for i, as := range tas {
		if i != asIndex {
			authsession = append(authsession, as)
			asmap[as.Token] = j
			j++
		}
	}

	return true, nil
}

func IsSessionTokenValid(username, clientip, token string) (bool, error) {
	if len(authsession) == 0 {
		return false, errors.New("No active session found")
	}
	asIndex := asmap[token]

	//// Check all the failing conditions
	if asIndex < 0 || asIndex > len(authsession) {
		return false, errors.New("Cannot find valid index in asmap")
	}

	if authsession[asIndex].Username != username {
		return false, errors.New("Username does not match in AuthSession")
	}

	if authsession[asIndex].Token != token {
		return false, errors.New("Token does not match in AuthSession")
	}

	if authsession[asIndex].ClientIp != clientip {
		return false, errors.New("ClientIP does not match in AuthSession")
	}

	thisActivityTime := time.Now().UTC().Unix()
	if thisActivityTime > authsession[asIndex].TokenHardExpiry {
		return false, errors.New("Token expired. New session Token required")
	}

	// TODO: v1.5
	/*
		if thisActivityTime > (authsession[asIndex].LastActivity + authsession[asIndex].IdleTimeOut) {
			return false
		}
	*/

	return true, nil
}

func IsAuthenticatedSession(w http.ResponseWriter, r *http.Request) bool {

	return true
}

func IsAuthenticatedAdminSession(w http.ResponseWriter, r *http.Request) bool {
	if !IsAuthenticatedSession(w, r) {
		return false
	}

	//// Implement code to verify it is admin session
	return true
}

func GetUseranmeByToken(token string) string {
	if token == "" {
		return ""
	}

	asIndex := asmap[token]
	//// Check all the failing conditions
	if asIndex < 0 || asIndex > len(authsession) {
		return ""
	}

	return authsession[asIndex].Username
}

//=================================================
// Login/Logff Handlers
//=================================================

func Login(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read HTTP Reuqest", err)
		return
	}

	var uacc Account
	uacc.LoadParams(body)

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

	boolUserExits, sacc := UserExists(uacc.Username)
	if !boolUserExits {
		HttpErrorReturn(w, r, "User does not exist", errors.New("User does not exist"))
		return
	}

	passwordMatched := CheckPasswordHash(uacc.Password, sacc.PasswordHash)
	if passwordMatched {
		clientip := "127.0.0.1" // TODO: Extract from request header
		sessionToken, err := CreateSession(uacc.Username, clientip)
		if err != nil {
			HttpErrorReturn(w, r, "Error authenticating user", errors.New("Failed to generate session Token"))
			return
		}
		w.Header().Set("NHR-TOKEN", sessionToken)
		fmt.Fprintln(w, "Login successful")
	} else {
		HttpErrorReturn(w, r, "Password did not match", errors.New("Password did not match"))
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HttpErrorReturn(w, r, "Failed to read HTTP Reuqest", err)
		return
	}

	postData := make(map[string]string)
	json.Unmarshal(body, &postData)

	postData["token"] = r.Header.Get("NHR-TOKEN")
	postData["clientip"] = "127.0.0.1"

	_, err = DestroySession(postData["username"], postData["clientip"], postData["token"])
	if err != nil {
		HttpErrorReturn(w, r, "Error destryoing session", err)
		return
	}

	HttpSuccessReturn(w, r, "Logout completed", 1)
}
