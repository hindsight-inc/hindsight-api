package topic

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Topic struct {
	gorm.Model
	ID string `json:"ID"`
	Title string `json:"title"`
	Content string `json:"content"`
	// Deadline Date
	// Author User
	// Options [Option]
	// PermissionView [User]
	// PermissionVote [User]
	// Cover Image
}

func (self *Topic) Response() gin.H {
	return gin.H{
		"title": self.Title,
		"content": self.Content,
	}
}

/*
	http://localhost:8080/topics
*/
func List(c *gin.Context) {
	topic1 := Topic{Title: "Mock Title 01", Content: "Mock content 01."}
	topic2 := Topic{Title: "Mock Title 02", Content: "Mock content 02."}
	topic3 := Topic{Title: "Mock Title 03", Content: "Mock content 03."}
	topics := []gin.H{
		topic1.Response(),
		topic2.Response(),
		topic3.Response(),
	}
	c.JSON(200, topics)
}

/*
	http://localhost:8080/topics/42
*/
func Detail(c *gin.Context) {
	id := c.Param("id")
	topic := Topic{ID: id, Title: "Mock Title 01", Content: "This is mock content 01.\n这是一个UTF8测试。"}
	c.JSON(200, topic)
}

/*
curl -v POST \
  http://localhost:8080/topics \
  -H 'content-type: application/json' \
  -d '{ "title": "Title 001", "content": "This is test contents." }'
*/
func Create(c *gin.Context) {
	var topic Topic
	c.BindJSON(&topic)
	c.JSON(200, topic)
}