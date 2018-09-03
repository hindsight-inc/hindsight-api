package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/*
curl -v -X POST \
  http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/login", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		if json.Username != "username001" || json.Password != "password001" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		} 
		
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	router.Run()
}