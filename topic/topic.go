package topic

import (
	"time"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"hindsight/user"
	"hindsight/database"
	"hindsight/error"
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
	Opinions	[]Opinion
	Votes		[]Vote
}

type CreateRequest struct {
	Title	string
	Content	string
	MilestoneDeadline time.Time `json:"milestone_deadline"`
	Opinions []OpinionRequest
}

func (self *Topic) DetailResponse() (int, gin.H) {
	db := database.GetDB()
	//	TODO: how to get gin.H from struct?
	if err := db.Model(self).Related(&self.Author, "Author").Error; err != nil {
		return error.Bad(error.DomainTopicResponse, error.ReasonDatabaseError, err.Error())
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
	}
}

/* Opinion */

type Opinion struct {
	gorm.Model
	Title		string
	TopicID		uint
	Author		user.User
	AuthorID	uint
}

type OpinionRequest struct {
	Title		string
}

/* Vote */

type Vote struct {
	gorm.Model
	Topic		Topic
	TopicID		uint
	Author		user.User
	AuthorID	uint
}