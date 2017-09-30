package models

import (
	"time"
)

// SSHKey struct
type SSHKey struct {
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

// BeforeInsert hooks
func (sshKey *SSHKey) BeforeInsert() {
	sshKey.CreatedUnix = time.Now().Unix()
	sshKey.UpdatedUnix = sshKey.CreatedUnix
}

// BeforeUpdate hooks
func (sshKey *SSHKey) BeforeUpdate() {
	sshKey.UpdatedUnix = time.Now().Unix()
}
