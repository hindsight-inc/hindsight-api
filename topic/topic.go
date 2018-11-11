package topic

import (
	"time"
	"github.com/jinzhu/gorm"
	"hindsight/user"
)

const kPageSize = "10"
const kTitleMin = 10
const kTitleMax = 1024

/* Topic */

type Topic struct {
	gorm.Model
	Title	string `binding:"required"`
	Content	string
	DeadlineStart	time.Time
	DeadlineEnd		time.Time `binding:"required"`
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
	DeadlineStart time.Time
}

/* Opinion */

type Opinion struct {
	gorm.Model
	Title		string
	TopicID		uint
	Author		user.User
	AuthorID	uint
}

/* Vote */

type Vote struct {
	gorm.Model
	Topic		Topic
	TopicID		uint
	Author		user.User
	AuthorID	uint
}