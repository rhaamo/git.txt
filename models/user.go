package models

import (
	"time"
)

type User struct {
	ID		int64	`xorm:"pk autoincr"`
	UserName	string	`xorm:"UNIQUE NOT NULL"`
	Password	string	`xorm:"NOT NULL"`
	Email		string	`xorm:"NOT NULL"`

	// Permissions
	IsAdmin		bool	`xorm:"DEFAULT 0"`
	IsActive	bool	`xorm:"DEFAULT 0"`

	Created		time.Time	`xorm:"-"`
	CreatedUnix	int64
	Updated		time.Time	`xorm:"-"`
	UpdatedUnix	int64

	// Relations
	// 	Gitxts
	// 	SshKeys
}

type SshKey struct {
	ID		int64	`xorm:"pk autoincr"`
	UserID		int64	`xorm:"INDEX NOT NULL"`

	Name		string	`xorm:"NOT NULL"`
	Fingerprint	string	`xorm:"NOT NULL"`
	Content		string	`xorm:"TEXT NOT NULL"`
	Mode		int	`xorm:"NOT NULL DEFAULT 2"`
	Type		int	`xorm:"NOT NULL DEFAULT 1"`

	Created		time.Time	`xorm:"-"`
	CreatedUnix	int64
	Updated		time.Time	`xorm:"-"`
	UpdatedUnix	int64

	// Relations
	// 	UserID
}