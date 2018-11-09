package topic

import (
	"github.com/jinzhu/gorm"
	"hindsight/user"
)

const kPageSize = "10"

/* Topic */

type Topic struct {
	gorm.Model
	Title	string
	Content	string
	// Deadline Date
	Author		user.User
	AuthorID	uint
	// PermissionView [User]
	// PermissionVote [User]
	// Cover Image
	Opinions	[]Opinion
	Votes		[]Vote
}

type TopicCreator struct {
	Title	string
	Content	string
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
