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
const kTitleMax = 140
const kContentMax = 65536
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
	CoverUID	string
	Opinions	tOpinions
	Votes		tVotes
}

type tTopics []Topic

/* Request */

type CreateRequest struct {
	Title		string
	Content		string
	MilestoneDeadline time.Time `json:"milestone_deadline"`
	CoverUID	string `json:"cover_uid"`
	Opinions	[]OpinionRequest
}

/* Response */

/// A simple response that doesn't replies on other database calls
func (self *Topic) CreateResponse() (int, gin.H) {
	return http.StatusOK, gin.H{
		"id": self.ID,
		"title": self.Title,
		"content": self.Content,
		"cover_uid": self.CoverUID,
		"milestone_deadline": self.MilestoneDeadline,
	}
}

func (self *Topic) Response() (int, gin.H) {
	db := database.Shared()
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
	code, hAuthor := self.Author.DetailResponse()
	if code != http.StatusOK {
		return code, hAuthor
	}
	code, hVotes := self.Votes.Response()
	if code != http.StatusOK {
		return code, hVotes[0]
	}
	return http.StatusOK, gin.H{
		"id": self.ID,
		"title": self.Title,
		"content": self.Content,
		"cover_uid": self.CoverUID,
		"milestone_deadline": self.MilestoneDeadline,
		"author": hAuthor,
		"opinions": self.Opinions.Response(),
		"votes": hVotes,
	}
}

func (self *Topic) DetailResponse() (int, gin.H) {
	db := database.Shared()
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
	code, hAuthor := self.Author.DetailResponse()
	if code != http.StatusOK {
		return code, hAuthor
	}
	code, hVotes := self.Votes.Response()
	if code != http.StatusOK {
		return code, hVotes[0]
	}
	return http.StatusOK, gin.H{
		"id": self.ID,
		"title": self.Title,
		"content": self.Content,
		"cover_uid": self.CoverUID,
		"milestone_deadline": self.MilestoneDeadline,
		"author": hAuthor,
		"opinions": self.Opinions.Response(),
		"votes": hVotes,
	}
}

func (self *tTopics) Response() (int, []gin.H) {
	var topics []gin.H
	for _, topic := range *self {
		code, h := topic.Response()
		if code != http.StatusOK {
			topics = nil
			return code, append(topics, h)
		}
		topics = append(topics, h)
	}
	return http.StatusOK, topics
}
