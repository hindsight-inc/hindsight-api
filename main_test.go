package main

import (
	"log"
	//"time"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"strconv"
	"testing"

	//"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"hindsight/auth"
	"hindsight/user"
	"hindsight/topic"
	"hindsight/error"
	"hindsight/config"
)

const kTestUserUsername = "test001"
const kTestUserPassword = "password123"
const kTestTopicTitle = "Script Test 测试"
const kTestTopicContent = "Test contents from script.\n测试内容"
const kSomething = "sth"


//	API Unit Tests

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


//	API Integration Tests

var UserID uint

func TestSetup(t *testing.T) {
	// no checking here as they just `panic`
	setupConfig()

	db := setupDB()
	defer db.Close()

	setupFacebook()

	middleware := setupAuth()
	assert.NotEmpty(t, middleware)
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
	u := user.User{Username: kTestUserUsername, Password: kTestUserPassword}
	b, _ := json.Marshal(u)

	//	permanently delete previous test users
	db.Unscoped().Where(user.User{Username: u.Username}).Delete(&user.User{})

	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	json.Unmarshal([]byte(w.Body.String()), &u)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, u.ID > 0)

	UserID = u.ID
}

func TestUserRegisterFailureDuplicated(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: kTestUserUsername, Password: kTestUserPassword}
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
	u := user.User{Username: kTestUserUsername, Password: kTestUserPassword}
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
	//assert.Equal(t, kTestUserUsername, u.Username)
}

func TestUserLoginFailureNonexistent(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: kTestUserUsername + kSomething, Password: kTestUserPassword}
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
	u := user.User{Username: kTestUserUsername, Password: kTestUserPassword + kSomething}
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
	assert.Contains(t, w.Body.String(), kTestUserUsername)
}

/*
curl -v GET \
  http://localhost:8080/user \
  -H 'Content-Type: application/json' \
  -H 'Authorization:Bearer xxx'
*/
func TestUserDetailSuccess(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/user", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	//	something like { "facebook_name": "N/A", "id": 2, "username": "test001" }
	var u user.User
	json.Unmarshal([]byte(w.Body.String()), &u)
	assert.Equal(t, u.Username, kTestUserUsername)
	assert.Contains(t, w.Body.String(), "N/A")	// <- TODO: this test shall break when we actually work on user detail API
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

var TopicID uint

/*
curl -v POST \
  http://localhost:8080/topics \
  -H 'content-type: application/json' \
  -H 'Authorization:Bearer xxx' \
  -d '{ "title": "Title Test 001", "content": "Test contents from script." }'
*/
func TestTopicCreateSuccess(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	t1 := topic.Topic{Title: kTestTopicTitle, Content: kTestTopicContent}
	b, _ := json.Marshal(t1)

	//	permanently delete previous test topics
	db.Unscoped().Where(topic.Topic{Title: t1.Title, Content: t1.Content}).Delete(&topic.Topic{})

	req, _ := http.NewRequest("POST", "/topics", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var t2 topic.Topic
	json.Unmarshal([]byte(w.Body.String()), &t2)
	assert.Equal(t, t2.Title, kTestTopicTitle)
	assert.Equal(t, t2.Content, kTestTopicContent)
	assert.True(t, t2.ID > 0)

	TopicID = t2.ID
}

func TestTopicCreateFailureEmptyTitle(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	t1 := topic.Topic{Content: kTestTopicContent}
	b, _ := json.Marshal(t1)

	req, _ := http.NewRequest("POST", "/topics", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, error.DomainTopicCreate, e.Domain)
	assert.Equal(t, error.ReasonNonexistentEntry, e.Reason)
}

/*
curl -v GET \
  http://localhost:8080/topics?offset=0&limit=5 \
  -H 'content-type: application/json' \
  -H 'Authorization:Bearer xxx'
*/
func TestTopicList1(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/topics?offset=0&limit=1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	topics := make([]topic.Topic, 0)
	json.Unmarshal([]byte(w.Body.String()), &topics)
	t0 := topics[0]
	assert.Equal(t, t0.ID, TopicID)
	assert.Equal(t, t0.Title, kTestTopicTitle)
	assert.Equal(t, t0.Content, kTestTopicContent)
}

/*
curl -v GET \
  http://localhost:8080/topics/xxx \
  -H 'content-type: application/json' \
  -H 'Authorization:Bearer xxx'
*/
func TestTopicDetail(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/topics/" + strconv.FormatUint(uint64(TopicID), 10), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var t0 topic.Topic
	json.Unmarshal([]byte(w.Body.String()), &t0)
	assert.Equal(t, t0.ID, TopicID)
	assert.Equal(t, t0.Title, kTestTopicTitle)
	assert.Equal(t, t0.Content, kTestTopicContent)
}

//	{{host}}/user/connect { "method": "facebook", "access_token": "xxx" }
func TestUserConnectFacebookSuccess(t *testing.T) {
	cfg := config.Shared()
	if cfg.Facebook_disable_test {
		return
	}

	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	r := auth.ConnectRequest{Method: "facebook", AccessToken: cfg.Facebook_access_token}
	b, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", "/user/connect", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var token auth.Token
	json.Unmarshal([]byte(w.Body.String()), &token)
	assert.NotEmpty(t, token.Expire)	// TODO: equal to or later than `now`
	assert.NotEmpty(t, token.Token)		// TODO: validate it's a correct JWT token
	Token = token.Token
}

func TestUserConnectFacebookFailureBadToken(t *testing.T) {
	cfg := config.Shared()
	if cfg.Facebook_disable_test {
		return
	}

	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	r := auth.ConnectRequest{Method: "facebook", AccessToken: kSomething}
	b, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", "/user/connect", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, error.DomainAuthJWT, e.Domain)
	assert.Equal(t, error.ReasonUnauthorized, e.Reason)
}

func TestUserConnectFailureBadMethod(t *testing.T) {
	cfg := config.Shared()
	if cfg.Facebook_disable_test {
		return
	}

	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	r := auth.ConnectRequest{Method: kSomething}
	b, _ := json.Marshal(r)

	req, _ := http.NewRequest("POST", "/user/connect", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, error.DomainAuthJWT, e.Domain)
	assert.Equal(t, error.ReasonUnauthorized, e.Reason)
}

func TestFacebookUserDetailSuccess(t *testing.T) {
	cfg := config.Shared()
	if cfg.Facebook_disable_test {
		return
	}

	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/user", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + Token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	//	something like { "facebook_name": "Leo Superarts", "id": 2, "username": "fb_2332682216774407_bffr0g3mvbaob6nqs880" }
	var u user.User
	json.Unmarshal([]byte(w.Body.String()), &u)
	assert.Contains(t, w.Body.String(), "fb_")	// <- TODO: create a model
}