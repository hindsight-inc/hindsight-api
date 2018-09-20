package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"testing"

	//"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"hindsight/user"
	"hindsight/error"
)

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
	u := user.User{Username: "test001", Password: "password123"}
	b, _ := json.Marshal(u)

	//	permanently delete previous test users
	db.Unscoped().Where(user.User{Username: u.Username}).Delete(&user.User{})

	//	register test user
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	json.Unmarshal([]byte(w.Body.String()), &u)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, u.ID > 0)

	//	do it again should fail
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/user/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var e error.APIError
	json.Unmarshal([]byte(w.Body.String()), &e)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, error.DomainUserRegister, e.Domain)
	assert.Equal(t, error.ReasonDuplicatedEntry, e.Reason)
}