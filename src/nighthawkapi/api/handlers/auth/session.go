package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"
)

type AuthSession struct {
	Username        string
	Role            string
	AdminUser       bool
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

// CreateSession creates new authenticated session information
// for a user and stores in memory
func CreateSession(username string, role string, clientip string) (string, error) {
	var as AuthSession
	as.Username = username
	as.Role = role
	as.ClientIp = clientip

	if strings.ToLower(role) == "admin" {
		as.AdminUser = true
	}

	as.NewSession()

	authsession = append(authsession, as)
	asmap[as.Token] = len(authsession) - 1

	return as.Token, nil
}

// DestroySession deletes authenticated session information
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

// IsAuthenticatedSession verifies that if request is authenticated request
// returns false with message if not a valid authenticated session
func IsAuthenticatedSession(w http.ResponseWriter, r *http.Request) (bool, string) {
	sessionData := make(map[string]string)

	sessionData["clientip"] = "127.0.0.1"            // setting default to 127.0.0.1. Must be overwritten using x-forward-for
	sessionData["token"] = r.Header.Get("NHR-TOKEN") // adding token informaiton

	//// session validation
	if len(authsession) == 0 {
		return false, "No active session found"
	}
	asIndex := asmap[sessionData["token"]]

	//// Check all the failing conditions
	if asIndex < 0 || asIndex > len(authsession) {
		return false, "Invalid token provided"
	}

	if authsession[asIndex].Token != sessionData["token"] {
		return false, "Token does not exist in authenticated session"
	}

	if authsession[asIndex].ClientIp != sessionData["clientip"] {
		return false, "Invalid Client IP for authenticated session"
	}

	thisActivityTime := time.Now().UTC().Unix()
	if thisActivityTime > authsession[asIndex].TokenHardExpiry {
		return false, "Token has expired"
	}

	// TODO: v1.5
	/*
		if thisActivityTime > (authsession[asIndex].LastActivity + authsession[asIndex].IdleTimeOut) {
			return false
		}
	*/

	return true, "Authenticated session"
}

// IsAuthenticatedAdminSession function first check if it is a valid authenticated session
// followed by authented user is member of admin group
func IsAuthenticatedAdminSession(w http.ResponseWriter, r *http.Request) (bool, string) {
	validSession, sessionMessage := IsAuthenticatedSession(w, r)
	if !validSession {
		return false, sessionMessage
	}

	//// Implement code to verify it is admin session
	sessionToken := r.Header.Get("NHR-TOKEN")

	if len(authsession) == 0 {
		return false, "No active authenticated session found"
	}
	asIndex := asmap[sessionToken]

	//// Check all the failing conditions
	if asIndex < 0 || asIndex > len(authsession) {
		return false, "Token is invalid or expired"
	}

	if authsession[asIndex].AdminUser {
		return true, "Authenticated admin user"
	}

	return false, "Authenticated user does not have admin role" //default return
}

// GetUsernameByToken function returns login username for authenticated token.
// This function returns empty string in error condition
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
