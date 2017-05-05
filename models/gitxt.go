package models

import (
	"time"
)

type Gitxt struct {
	ID		int64	`xorm:"pk autoincr"`
	Hash		string	`xorm:"UNIQUE NOT NULL"`
	UserID		int64	`xorm:"INDEX NOT NULL"`
	Description	string	`xorm:"TEXT"`

	// Permissions
	IsPrivate	bool	`xorm:"DEFAULT 0"`

	Created		time.Time	`xorm:"-"`
	CreatedUnix	int64
	Updated		time.Time	`xorm:"-"`
	UpdatedUnix	int64

	// Relations
	// 	UserID
}