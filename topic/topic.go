package topic

import (
	"github.com/jinzhu/gorm"
	"hindsight/user"
)

const kPageSize = "10"

type Topic struct {
	gorm.Model
	Title	string `json:"title"`
	Content	string `json:"content"`
	// Deadline Date
	Author		user.User
	AuthorID	uint
	// Options [Option]
	// PermissionView [User]
	// PermissionVote [User]
	// Cover Image
}

type TopicCreator struct {
	Title	string `json:"title"`
	Content	string `json:"content"`
}