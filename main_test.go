package main

import (
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"testing"

	//"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"hindsight/user"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestUserRegister(t *testing.T) {
	db := setupDB()
	defer db.Close()
	router := setupRouter()

	w := httptest.NewRecorder()
	u := user.User{Username: "test001", Password: "password123"}
	b, _ := json.Marshal(u)

	//	permanently delete previous test users
	db.Unscoped().Where(user.User{Username: u.Username}).Delete(&user.User{})

	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}