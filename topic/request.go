package topic

import (
	//"log"
	"time"
	//"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"hindsight/database"
	"hindsight/user"
	"hindsight/herror"
)

func List(c *gin.Context) {
	offset := c.DefaultQuery("offset", "0")
	limit := c.DefaultQuery("limit", kPageSize)
	db := database.GetDB()
	var topics []Topic
	if err := db.Order("updated_at desc, created_at desc").Offset(offset).Limit(limit).Find(&topics).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, topics)
}

func Detail(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")
	var topic Topic
	if err := db.First(&topic, id).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
		return
	}
	c.JSON(topic.DetailResponse())
}

func VoteOpinion(c *gin.Context) {
	db := database.GetDB()
	tid := c.Param("id")
	oid := c.Param("oid")

	//	Get topic
	var topic Topic
	if err := db.First(&topic, tid).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
		return
	}

	//	Get opinion
	var opinion Opinion
	if err := db.First(&opinion, oid).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
		return
	}
	if opinion.TopicID != topic.ID {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonInvalidEntry, "Opinion does not belong to this topic"))
		return
	}

	//	Get author
	u := user.Current(c)
	if u == nil {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonNonexistentEntry, "User not found"))
		return
	}

	//	Create vote
	vote := Vote{
		AuthorID: u.ID,
		TopicID: topic.ID,
	}
	//	check if user already voted
	if err := db.First(&vote).Error; err != nil && vote.ID > 0 {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
		return
	}
	//	TODO: re-vote logic will be decided after UI is done
	if vote.ID > 0 {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDuplicatedEntry, "Cannot re-vote"))
		return
	}
	vote.OpinionID = opinion.ID
	if err := db.Create(&vote).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
		return
	}

	//	TODO: return something useful, like updated votes count for the topic
	c.JSON(http.StatusOK, vote.Response())
}

func Create(c *gin.Context) {
	u := user.Current(c)
	if u == nil {
		//	Shouldn't reach here unless user has been deleted but active token is not
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonNonexistentEntry, "User not found"))
		return
	}

	db := database.GetDB()
	var request CreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonInvalidJSON, err.Error()))
		return
	}
	//	TODO: we are not doing null check for required fields because I don't know how to do it in golang.
	//	API still fails due to ReasonInvalidEntry instead of ReasonNonexistentEntry, but the message will be misleading.

	//	unicode/utf8.RuneCountInString counts characters correctly, 
	//	but 8 Chinese characters can make a valid title.
	//	e.g. 川普年底会倒台吗 (Will Trump be impeached by the end of this year)
	if len(request.Title) < kTitleMin {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonInvalidEntry, "title is too short"))
		return
	}
	if len(request.Title) > kTitleMax {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonInvalidEntry, "title is too long"))
		return
	}
	if request.MilestoneDeadline.Before(time.Now().Add(kDeadlineThreshold)) {
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonNonexistentEntry, "milestone_deadline is not late enough"))
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
		c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
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
		opinion := Opinion{
			Title: o.Title,
			TopicID: topic.ID,
			AuthorID: u.ID,
		}
		if err := db.Create(&opinion).Error; err != nil {
			c.JSON(herror.Bad(herror.DomainTopicCreate, herror.ReasonDatabaseError, err.Error()))
			//	TODO: revert topic creation
			return
		}
	}
	c.JSON(http.StatusOK, topic.Response())
}