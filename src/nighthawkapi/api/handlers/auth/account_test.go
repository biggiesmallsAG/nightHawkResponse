package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

<<<<<<< HEAD
func TestCreateNewUser(t *testing.T) {
	var acc Account
	acc.Firstname = "roshan"
	acc.Lastname = "maskey"
	acc.Username = "admin"
	acc.Password = "admin"
	acc.Role = "admin"
=======
//// Test new account creation
func TestCreateUser(t *testing.T) {

	// Step 1: Authenticate as admin user
	loginData := make(map[string]string)
	loginData["username"] = "admin"
	loginData["password"] = "admin"
	jLoginData, _ := json.Marshal(loginData)

	rr := loginUser(t, jLoginData)
	sessionToken := rr.Header().Get("NHR-TOKEN")

	// Step 2: Created new account as authenticated user
	acc := Account{
		Firstname: "roshan",
		Lastname:  "maskey",
		Email:     "roshanmaskey@nighthawk.local",
		Username:  "roshan8",
		Password:  "roshan",
		Role:      "admin",
	}
>>>>>>> 5af04688964ecae66fb4462369b904a79f6d5af8

	jd, _ := json.Marshal(acc)
	createNewUser(t, sessionToken, jd)

	// Step 3: Logout
	logoutData := make(map[string]string)
	logoutData["username"] = "admin"
	jlogoutData, _ := json.Marshal(logoutData)
	logout(t, sessionToken, jlogoutData)
}

func createNewUser(t *testing.T, sessionToken string, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/admin/user/create", bytes.NewReader(data))
	req.Header.Set("NHR-TOKEN", sessionToken)
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

//// Test useraccount deletion
func TestDeleteUser(t *testing.T) {

	// Step 1: Authenticate as admin user
	loginData := make(map[string]string)
	loginData["username"] = "admin"
	loginData["password"] = "admin"
	jLoginData, _ := json.Marshal(loginData)

	rr := loginUser(t, jLoginData)
	sessionToken := rr.Header().Get("NHR-TOKEN")

	// Step 2: Delete user by username
	userinfo := make(map[string]string)
	userinfo["username"] = "roshan8"
	data, _ := json.Marshal(userinfo)

	deleteUser(t, sessionToken, data)

	// Step 3: Logout
	logoutData := make(map[string]string)
	logoutData["username"] = loginData["username"]
	jlogoutData, _ := json.Marshal(logoutData)
	logout(t, sessionToken, jlogoutData)
}

func deleteUser(t *testing.T, token string, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/admin/user/delete", bytes.NewReader(data))
	req.Header.Set("NHR-TOKEN", token)

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

// Test SetPassword
func TestSetPassword(t *testing.T) {
	// Step 1: Authenticate as admin user
	loginData := make(map[string]string)
	loginData["username"] = "admin"
	loginData["password"] = "admin"
	jLoginData, _ := json.Marshal(loginData)

	rr := loginUser(t, jLoginData)
	sessionToken := rr.Header().Get("NHR-TOKEN")

	// Step 2: Set password
	userinfo := make(map[string]string)
	userinfo["username"] = "roshan8"
	userinfo["new_password"] = "password123"
	data, _ := json.Marshal(userinfo)

	setPassword(t, sessionToken, data)

	// Step 3: Logout
	logoutData := make(map[string]string)
	logoutData["username"] = loginData["username"]
	jlogoutData, _ := json.Marshal(logoutData)
	logout(t, sessionToken, jlogoutData)
}

func setPassword(t *testing.T, token string, data []byte) {
	req, err := http.NewRequest("POST", "/api/v1/admin/password/set", bytes.NewReader(data))
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

func loginUser(t *testing.T, jd []byte) *httptest.ResponseRecorder {
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

//// Test changing password
func TestChangePassword(t *testing.T) {

	// Step 1: Authenticate as user
	loginData := make(map[string]string)
	loginData["username"] = "roshan4"
	loginData["password"] = "roshan"
	jLoginData, _ := json.Marshal(loginData)

	rr := loginUser(t, jLoginData)
	sessionToken := rr.Header().Get("NHR-TOKEN")

	// Step 2: Change current password
	passwdData := make(map[string]string)
	passwdData["password"] = loginData["password"]
	passwdData["new_password"] = "roshan123"
	data, _ := json.Marshal(passwdData)

	changePassword(t, sessionToken, data)

	// Step 3: Logout
	logoutData := make(map[string]string)
	logoutData["username"] = loginData["username"]
	jlogoutData, _ := json.Marshal(logoutData)
	logout(t, sessionToken, jlogoutData)

	time.Sleep(1 * time.Second)

	// Step 4: Login Using new password
	loginData["password"] = passwdData["new_password"]
	data, _ = json.Marshal(loginData)
	loginUser(t, data)
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

	rr := loginUser(t, jdata)
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
