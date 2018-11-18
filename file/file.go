package file

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/jinzhu/gorm"

	"hindsight/herror"
	"hindsight/database"
)

const ImagePath = "./public/upload/image/"

type Image struct {
	gorm.Model
	Filename string
	ContentType string
	Title string
	UID string
}

func (self *Image) Response() (int, gin.H) {
	return http.StatusOK, gin.H{
		"uid": self.UID,
	}
}

func UploadImage(c *gin.Context) {
	title := c.PostForm("title")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(herror.Bad(herror.DomainUploadImage, herror.ReasonInvalidEntry, err.Error()))
		return
	}

	mime := file.Header.Get("Content-Type")
	if !strings.Contains(mime, "image") {
		c.JSON(herror.Bad(herror.DomainUploadImage, herror.ReasonInvalidEntry, "Invalid content-type: " + mime))
		return
	}

	uid := xid.New().String()
	//	xid needs to be truly unique
	//	TODO: check if file is image
	//	TODO: check file size

	if err := c.SaveUploadedFile(file, ImagePath + uid); err != nil {
		c.JSON(herror.Bad(herror.DomainUploadImage, herror.ReasonOperationFailure, err.Error()))
		return
	}

	image := Image{
		Filename: file.Filename,
		ContentType: mime,
		Title: title,
		UID: uid,
	}

	db := database.Shared()
	if err := db.Create(&image).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainUploadImage, herror.ReasonDatabaseError, err.Error()))
		return
	}

	c.JSON(image.Response())
}

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