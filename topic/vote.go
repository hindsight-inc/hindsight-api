package topic

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"hindsight/user"
	"hindsight/database"
	"hindsight/herror"
)

type Vote struct {
	gorm.Model
	Topic		Topic
	TopicID		uint
	Opinion		Opinion
	OpinionID	uint
	Author		user.User
	AuthorID	uint
}

type tVotes []Vote

/* Response */

func (self *Vote) Response() (int, gin.H) {
	db := database.Shared()
	if err := db.Model(self).Related(&self.Author, "Author").Error; err != nil {
		return herror.Bad(herror.DomainTopicResponse, herror.ReasonDatabaseError, err.Error())
	}
	code, h := self.Author.DetailResponse()
	if code != http.StatusOK {
		return code, h
	}
	return http.StatusOK, gin.H{
		"id": self.ID,
		"topic_id": self.TopicID,
		"opinion_id": self.OpinionID,
		"author_id": self.AuthorID,
	}
}

func (self *Vote) DetailResponse() (int, gin.H) {
	db := database.Shared()
	if err := db.Model(self).Related(&self.Author, "Author").Error; err != nil {
		return herror.Bad(herror.DomainTopicResponse, herror.ReasonDatabaseError, err.Error())
	}
	code, h := self.Author.DetailResponse()
	if code != http.StatusOK {
		return code, h
	}
	return http.StatusOK, gin.H{
		"id": self.ID,
		"topic_id": self.TopicID,
		"opinion_id": self.OpinionID,
		"author_id": self.AuthorID,
		"author": h,
	}
}

func (self *tVotes) Response() (int, []gin.H) {
	var votes []gin.H
	for _, vote := range *self {
		code, h := vote.Response()
		if code != http.StatusOK {
			votes = nil
			return code, append(votes, h)
		}
		votes = append(votes, h)
	}
	return http.StatusOK, votes
}

func (self *tVotes) DetailResponse() (int, []gin.H) {
	var votes []gin.H
	for _, vote := range *self {
		code, h := vote.DetailResponse()
		if code != http.StatusOK {
			votes = nil
			return code, append(votes, h)
		}
		votes = append(votes, h)
	}
	return http.StatusOK, votes
}