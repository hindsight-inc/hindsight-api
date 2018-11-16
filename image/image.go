package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/", "../public/upload/image")
	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")

		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		mime := file.Header.Get("Content-Type")
		if !strings.Contains(mime, "image") {
			c.String(http.StatusBadRequest, fmt.Sprintf("not image: %s", mime))
			return
		}
		//	TODO: check if file is image
		//	TODO: check file size
		if err := c.SaveUploadedFile(file, "../public/upload/image/" + file.Filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s. type: %s", file.Filename, name, email, file.Header.Get("Content-Type")))
	})
	router.Run(":8080")
}