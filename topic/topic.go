package topic

import (
	"github.com/jinzhu/gorm"
)

const kPageSize = "10"

type Topic struct {
	gorm.Model
	Title string `json:"title"`
	Content string `json:"content"`
	// Deadline Date
	// Author User
	// Options [Option]
	// PermissionView [User]
	// PermissionVote [User]
	// Cover Image
}