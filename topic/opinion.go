package topic

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"hindsight/user"
)

type Opinion struct {
	gorm.Model
	Title		string
	TopicID		uint
	Author		user.User
	AuthorID	uint
	VoteCount	uint `gorm:"-"`
}

type tOpinions []Opinion

/* Request */

type OpinionRequest struct {
	Title		string
}

/* Response */

func (self *Opinion) CountResponse() gin.H {
	return gin.H{
		"id": self.ID,
		"title": self.Title,
		"vote_count": self.VoteCount,
	}
}

func (self *tOpinions) CountResponse() []gin.H {
	var opinions []gin.H
	for _, opinion := range *self {
		opinions = append(opinions, opinion.CountResponse())
	}
	return opinions
}