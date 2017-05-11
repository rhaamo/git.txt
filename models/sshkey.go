package models

import (
	"time"
)

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

func (ssh_key *SshKey) BeforeInsert() {
	ssh_key.CreatedUnix = time.Now().Unix()
	ssh_key.UpdatedUnix = ssh_key.CreatedUnix
}

func (ssh_key *SshKey) BeforeUpdate() {
	ssh_key.UpdatedUnix = time.Now().Unix()
}
