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
}

type tOpinions []Opinion

/* Request */

type OpinionRequest struct {
	Title		string
}

/* Response */

func (self *Opinion) Response() gin.H {
	return gin.H{"title": self.Title}
}

func (self *tOpinions) Response() []gin.H {
	var opinions []gin.H
	for _, opinion := range *self {
		opinions = append(opinions, opinion.Response())
	}
	return opinions
}
