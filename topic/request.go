package topic

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"hindsight/database"
	"hindsight/user"
	"hindsight/error"
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
	u := user.Current(c)
	if u == nil {
		//	Shouldn't reach here unless user has been deleted but active token is not
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonNonexistentEntry, "User not found"))
		return
	}

	db := database.GetDB()
	var request CreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//	unicode/utf8.RuneCountInString counts characters correctly, 
	//	but 8 Chinese characters can make a valid title.
	//	e.g. 川普年底会倒台吗 (Will Trump be impeached by the end of this year)
	if len(request.Title) < kTitleMin {
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonInvalidEntry, "Title is too short"))
		return
	}
	if len(request.Title) > kTitleMax {
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonInvalidEntry, "Title is too long"))
		return
	}
	topic := Topic{
		Title: request.Title,
		Content: request.Content,
		AuthorID: u.ID,
		DeadlineStart: time.Unix(3600, 0),
		DeadlineEnd: time.Unix(3600, 0),
	}
	db.Create(&topic)
	c.JSON(http.StatusOK, topic)
}