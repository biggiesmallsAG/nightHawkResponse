package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNeUser(t *testing.T) {
	var acc Account
	acc.Firstname = "roshan"
	acc.Lastname = "maskey"
	acc.Username = "roshan"
	acc.Password = "roshan"
	acc.Role = "admin"

	jd, _ := json.Marshal(acc)
	createNewUser(t, jd)

}

func createNewUser(t *testing.T, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/auth/user/create", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateNewUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println(rr.Body.String())
}

func TestLogin(t *testing.T) {
	var acc1 Account
	acc1.Username = "admin"
	acc1.Password = "admin"

	jd, _ := json.Marshal(acc1)

	// Test if token is valid
	rr := nighthawkLogin(t, jd)
	sessionToken := rr.Header().Get("Nhr-Token")
	fmt.Println("NHR-TOKEN", sessionToken)

	// Testing for valid token
	tokenTest("admin", "127.0.0.1", sessionToken)

	// Testing for invalid token
	tokenTest("admin1", "127.0.0.1", sessionToken)
	tokenTest("admin", "127.0.0.2", sessionToken)
	tokenTest("admin", "127.0.0.1", sessionToken+"a")
}

func TestLogin2(t *testing.T) {
	acc := Account{Username: "admin", Password: "admin"}
	jd, _ := json.Marshal(acc)

	// Login1
	rr := nighthawkLogin(t, jd)
	sessionToken1 := rr.Header().Get("Nhr-Token")
	fmt.Println("NHR-TOKEN", sessionToken1)

	// Login2
	// Test if token is valid
	rr = nighthawkLogin(t, jd)
	sessionToken2 := rr.Header().Get("Nhr-Token")
	fmt.Println("NHR-TOKEN", sessionToken2)

	rr = nighthawkLogin(t, jd)
	sessionToken3 := rr.Header().Get("Nhr-Token")
	fmt.Println("NHR-TOKEN", sessionToken3)

	rr = nighthawkLogin(t, jd)
	sessionToken4 := rr.Header().Get("Nhr-Token")
	fmt.Println("NHR-TOKEN", sessionToken4)

	fmt.Println("Token test after creation")
	tokenTest("admin", "127.0.0.1", sessionToken4)
	tokenTest("admin", "127.0.0.1", sessionToken3)
	tokenTest("admin", "127.0.0.1", sessionToken2)
	tokenTest("admin", "127.0.0.1", sessionToken1)

	DestroySession("admin", "127.0.0.1", sessionToken2)

	fmt.Println("Token test after deletion")
	tokenTest("admin", "127.0.0.1", sessionToken4)
	tokenTest("admin", "127.0.0.1", sessionToken3)
	tokenTest("admin", "127.0.0.1", sessionToken2)
	tokenTest("admin", "127.0.0.1", sessionToken1)

	fmt.Println(asmap)

	fmt.Println(authsession)

}

func nighthawkLogin(t *testing.T, jd []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(jd))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println(rr.Body.String())
	return rr
}

func tokenTest(username, clientip, sessionToken string) {
	tokenStatus, err := IsSessionTokenValid(username, clientip, sessionToken)
	if tokenStatus {
		fmt.Println("Valid Token")
	} else {
		fmt.Println("Invalid Token ", err.Error())
	}
}

func TestChangePassword(t *testing.T) {

	passwdData := make(map[string]string)
	passwdData["username"] = "admin2"
	passwdData["password"] = "admin2"

	jpasswdData, _ := json.Marshal(passwdData)

	var res map[string]string
	createNewUser(t, jpasswdData)

	rr := nighthawkLogin(t, jpasswdData)
	json.Unmarshal([]byte(rr.Body.String()), &res)

	if res["response"] == "failed" {
		fmt.Println("Login failed")
		return
	}

	sessionToken := rr.Header().Get("Nhr-Token")
	fmt.Println(sessionToken)

	passwdData["new_password"] = "maskey"
	jpasswdData, _ = json.Marshal(passwdData)
	fmt.Println(string(jpasswdData))

	changePassword(t, sessionToken, jpasswdData)
}

func changePassword(t *testing.T, token string, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/auth/password/change", bytes.NewReader(data))
	req.Header.Set("NHR-TOKEN", token)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ChangePassword)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println(rr.Body.String())
}

func TestDeleteUser(t *testing.T) {

	sessionToken := ""
	postData := make(map[string]string)
	postData["delete_account"] = "admin2"

	jd, _ := json.Marshal(postData)
	//fmt.Println(string(jd))
	deleteUser(t, sessionToken, jd)
}

func deleteUser(t *testing.T, token string, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/auth/user/delete", bytes.NewReader(data))
	//req.Header.Set("NHR-TOKEN", token)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println(rr.Body.String())
}

func TestSetPassword(t *testing.T) {
	postData := make(map[string]string)
	postData["username"] = "admin"
	postData["password"] = "admin"

	jsonPostData, _ := json.Marshal(postData)
	rr := nighthawkLogin(t, jsonPostData)

	res := make(map[string]string)
	json.Unmarshal([]byte(rr.Body.String()), &res)

	if res["result"] == "failed" || res["response"] == "failed" {
		fmt.Println("Login failed for ", postData["username"])
		return
	}

	sessionToken := rr.Header().Get("NHR-TOKEN")
	postData["username"] = "admin1"
	postData["new_password"] = "maskey2"
	jsonPostData, _ = json.Marshal(postData)

	setPassword(t, sessionToken, jsonPostData)
}

func setPassword(t *testing.T, token string, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/auth/password/set", bytes.NewReader(data))
	req.Header.Set("NHR-TOKEN", token)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SetPassword)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println(rr.Body.String())
}

func logout(t *testing.T, token string, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/auth/logout", bytes.NewReader(data))
	req.Header.Set("NHR-TOKEN", token)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Logout)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fmt.Println(rr.Body.String())
}

func TestLoginLogout(t *testing.T) {
	data := make(map[string]string)
	data["username"] = "admin1"
	data["password"] = "maskey2"
	jdata, _ := json.Marshal(data)

	rr := nighthawkLogin(t, jdata)
	res := make(map[string]string)
	json.Unmarshal([]byte(rr.Body.String()), &res)

	if res["response"] == "failed" {
		fmt.Println("Login failed: ", data["username"])
		return
	}

	sessionToken := rr.Header().Get("NHR-TOKEN")
	data["username"] = "admin1"
	data["password"] = ""
	jdata, _ = json.Marshal(data)

	logout(t, sessionToken, jdata)
}
