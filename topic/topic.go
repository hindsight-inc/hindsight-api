package topic

import (
	"time"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"hindsight/user"
	"hindsight/database"
	"hindsight/herror"
)

const kPageSize = "10"
const kTitleMin = 10
const kTitleMax = 1024
const kDeadlineThreshold = time.Hour * 1
const kDefaultTitle0 = "Agree"
const kDefaultTitle1 = "Disagree"

/* Topic */

type Topic struct {
	gorm.Model
	Title	string `binding:"required"`
	Content	string
	/*
	 We only have one milestone so far.
	 Potentially, we will have more milestones,
	 e.g. milestoneInvite, milestoneVote, milestoneEnd, etc.
	 */
	MilestoneDeadline time.Time `binding:"required"`
	Author		user.User
	AuthorID	uint
	// PermissionView [User]
	// PermissionVote [User]
	// Cover Image
	Opinions	tOpinions
	Votes		tVotes
}

/* Request */

type CreateRequest struct {
	Title	string
	Content	string
	MilestoneDeadline time.Time `json:"milestone_deadline"`
	Opinions []OpinionRequest
}

/* Response */

func (self *Topic) Response() gin.H {
	return gin.H{
		"id": self.ID,
		"title": self.Title,
		"content": self.Content,
		"milestone_deadline": self.MilestoneDeadline,
	}
}

func (self *Topic) DetailResponse() (int, gin.H) {
	db := database.GetDB()
	//	TODO: how to get gin.H from struct?
	if err := db.Model(self).Related(&self.Author, "Author").Error; err != nil {
		return herror.Bad(herror.DomainTopicResponse, herror.ReasonDatabaseError, err.Error())
	}
	if err := db.Model(self).Related(&self.Opinions, "Opinions").Error; err != nil {
		return herror.Bad(herror.DomainTopicResponse, herror.ReasonDatabaseError, err.Error())
	}
	if err := db.Model(self).Related(&self.Votes, "Votes").Error; err != nil {
		return herror.Bad(herror.DomainTopicResponse, herror.ReasonDatabaseError, err.Error())
	}
	code, h := self.Author.DetailResponse()
	if code != http.StatusOK {
		return code, h
	}
	return http.StatusOK, gin.H{
		"id": self.ID,
		"title": self.Title,
		"content": self.Content,
		"milestone_deadline": self.MilestoneDeadline,
		"author": h,
		"opinions": self.Opinions.Response(),
		"votes": self.Votes.Response(),
	}
}

/* Vote */

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

func (self *Vote) Response() gin.H {
	return gin.H{
		"id": self.ID,
		"topic_id": self.TopicID,
		"opinion_id": self.OpinionID,
		"author_id": self.AuthorID,
	}
}

func (self *tVotes) Response() []gin.H {
	var votes []gin.H
	for _, vote := range *self {
		votes = append(votes, vote.Response())
	}
	return votes
}