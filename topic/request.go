package topic

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"hindsight/database"
)

func (self *Topic) Response() gin.H {
	return gin.H{
		"title": self.Title,
		"content": self.Content,
	}
}

func List(c *gin.Context) {
	offset := c.DefaultQuery("offset", "0")
	limit := c.DefaultQuery("limit", kPageSize)
	db := database.GetDB()
	var topics []Topic
	db.Order("updated_at desc, created_at desc").Offset(offset).Limit(limit).Find(&topics)
	c.JSON(200, topics)
}

func Detail(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	var topic Topic
	db.First(&topic, id)
	c.JSON(200, topic)
}

func Create(c *gin.Context) {
	db := database.GetDB()
	var topic Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topic = Topic{Title: topic.Title, Content: topic.Content}
	db.Create(&topic)
	c.JSON(http.StatusOK, topic)
}