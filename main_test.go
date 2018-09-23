package main

import (
	"log"
	//"time"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"testing"

	//"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"hindsight/auth"
	"hindsight/user"
	"hindsight/error"
)

const kTestUsername = "test001"
const kTestPassword = "password123"
const kSomething = "sth"

/*
http://localhost:8080/ping
*/
func TestPingRoute(t *testing.T) {
	log.Println("Test: ping/pong")
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

/*
curl -v POST \
  http://localhost:8080/user/register \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func TestUserRegister(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: kTestUsername, Password: kTestPassword}
	b, _ := json.Marshal(u)

	//	permanently delete previous test users
	db.Unscoped().Where(user.User{Username: u.Username}).Delete(&user.User{})

	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	json.Unmarshal([]byte(w.Body.String()), &u)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, u.ID > 0)
}

func TestUserRegisterFailureDuplicated(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: kTestUsername, Password: kTestPassword}
	b, _ := json.Marshal(u)

	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, error.DomainUserRegister, e.Domain)
	assert.Equal(t, error.ReasonDuplicatedEntry, e.Reason)
}

var Token string

/*
curl -v POST \
  http://localhost:8080/user/login \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func TestUserLoginSuccess(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: kTestUsername, Password: kTestPassword}
	b, _ := json.Marshal(u)

	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	//	authMiddleware.LoginHandler
	var token auth.Token
	json.Unmarshal([]byte(w.Body.String()), &token)
	assert.NotEmpty(t, token.Expire)	// TODO: equal to or later than `now`
	assert.NotEmpty(t, token.Token)		// TODO: validate it's a correct JWT token
	Token = token.Token
	//	user.UserLogin
	//json.Unmarshal([]byte(w.Body.String()), &u)
	//assert.True(t, u.ID > 0)
	//assert.Equal(t, kTestUsername, u.Username)
}

func TestUserLoginFailureNonexistent(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: kTestUsername + kSomething, Password: kTestPassword}
	b, _ := json.Marshal(u)

	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	//	authMiddleware.LoginHandler
	assert.Equal(t, error.DomainAuthJWT, e.Domain)
	assert.Equal(t, error.ReasonUnauthorized, e.Reason)
	//	user.UserLogin
	//assert.Equal(t, error.DomainUserLogin, e.Domain)
	//assert.Equal(t, error.ReasonNonexistentEntry, e.Reason)
}

func TestUserLoginFailureMismatch(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: kTestUsername, Password: kTestPassword + kSomething}
	b, _ := json.Marshal(u)

	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, error.DomainAuthJWT, e.Domain)
	assert.Equal(t, error.ReasonUnauthorized, e.Reason)
}

/*
curl -v GET \
  http://localhost:8080/token/ping \
  -H 'Content-Type: application/json' \
  -H 'Authorization:Bearer xxx'
*/
func TestUserTokenPingSuccess(t *testing.T) {
	//	TODO: figure out why test would fail without `setupDB`
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/token/ping", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	//	something like {"claim_id":"test001","message":"pong","username":"test001"}
	assert.Contains(t, w.Body.String(), "pong")
	assert.Contains(t, w.Body.String(), kTestUsername)
}

func TestUserTokenPingFailureUnauthorized(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/token/ping", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token + kSomething)
	router.ServeHTTP(w, req)

	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, error.DomainAuthJWT, e.Domain)
	assert.Equal(t, error.ReasonUnauthorized, e.Reason)
}

/*
curl -v GET \
  http://localhost:8080/token/refresh \
  -H 'content-type: application/json' \
  -H 'Authorization:Bearer xxx'
*/
func TestUserTokenRefreshSuccess(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/token/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var token auth.Token
	json.Unmarshal([]byte(w.Body.String()), &token)
	assert.NotEmpty(t, token.Expire)	// TODO: equal to or later than `now`
	assert.NotEmpty(t, token.Token)		// TODO: validate it's a correct JWT token
	Token = token.Token
}

func TestUserTokenRefreshFailureUnauthorized(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/token/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token + kSomething)
	router.ServeHTTP(w, req)

	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, error.DomainAuthJWT, e.Domain)
	assert.Equal(t, error.ReasonUnauthorized, e.Reason)
}