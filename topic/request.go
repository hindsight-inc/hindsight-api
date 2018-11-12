package topic

import (
	"log"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"hindsight/database"
	"hindsight/user"
	"hindsight/error"
)

func List(c *gin.Context) {
	offset := c.DefaultQuery("offset", "0")
	limit := c.DefaultQuery("limit", kPageSize)
	db := database.GetDB()
	var topics []Topic
	db.Order("updated_at desc, created_at desc").Offset(offset).Limit(limit).Find(&topics)
	c.JSON(http.StatusOK, topics)
}

func Detail(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	var topic Topic
	db.First(&topic, id)
	c.JSON(topic.DetailResponse())
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
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonInvalidJSON, err.Error()))
		return
	}
	//	TODO: we are not doing null check for required fields because I don't know how to do it in golang.
	//	API still fails due to ReasonInvalidEntry instead of ReasonNonexistentEntry, but the message will be misleading.

	//	unicode/utf8.RuneCountInString counts characters correctly, 
	//	but 8 Chinese characters can make a valid title.
	//	e.g. 川普年底会倒台吗 (Will Trump be impeached by the end of this year)
	if len(request.Title) < kTitleMin {
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonInvalidEntry, "title is too short"))
		return
	}
	if len(request.Title) > kTitleMax {
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonInvalidEntry, "title is too long"))
		return
	}
	if request.MilestoneDeadline.Before(time.Now().Add(kDeadlineThreshold)) {
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonNonexistentEntry, "milestone_deadline is not late enough"))
		return
	}

	//	Create topic
	topic := Topic{
		Title: request.Title,
		Content: request.Content,
		AuthorID: u.ID,
		MilestoneDeadline: request.MilestoneDeadline,
	}
	if err := db.Create(&topic).Error; err != nil {
		c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonDatabaseError, err.Error()))
		return
	}

	//	Create opinions
	opinions := request.Opinions
	if len(opinions) == 0 {
		//	Add default opinions; should this logic be handled on front-end side?
		opinions = append(opinions, OpinionRequest{
			Title: kDefaultTitle0,
		}, OpinionRequest{
			Title: kDefaultTitle1,
		})
	}

	for _, o := range opinions {
		log.Println(o)
		opinion := Opinion{
			Title: o.Title,
			TopicID: topic.ID,
			AuthorID: u.ID,
		}
		if err := db.Create(&opinion).Error; err != nil {
			c.JSON(error.Bad(error.DomainTopicCreate, error.ReasonDatabaseError, err.Error()))
			//	TODO: revert topic creation
			return
		}
	}
	c.JSON(http.StatusOK, topic.Response())
}