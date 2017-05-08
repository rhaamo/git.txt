package models

import (
	"time"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
)

type Gitxt struct {
	ID		int64	`xorm:"pk autoincr"`
	Hash		string	`xorm:"UNIQUE NOT NULL"`
	UserID		int64	`xorm:"INDEX"`
	Anonymous	bool
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

// IsHashUsed checks if given hash exist,
func IsHashUsed(uid int64, hash string) (bool, error) {
	if len(hash) == 0 {
		return false, nil
	}
	return x.Get(&Gitxt{Hash: hash})
}

// Create a new gitxt
func CreateGitxt(g *Gitxt) (err error) {
	isExist, err := IsHashUsed(0, g.Hash)
	if err != nil {
		return err
	} else if isExist {
		return ErrHashAlreadyExist{g.Hash}
	}

	sess := x.NewSession()
	defer sessionRelease(sess)
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Insert(g); err != nil {
		return err
	}

	return sess.Commit()
}

// GetRepositoryByName returns the repository by given name under user if exists.
func GetRepositoryByName(user string, name string) (*Gitxt, error) {
	// First get user
	u, err := GetUserByName(user)
	if err != nil {
		return nil, err
	}

	repo := &Gitxt{
		UserID:   u.ID,
		Hash:	  name,
	}
	has, err := x.Get(repo)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.RepoNotExist{0, u.ID, name}
	}
	return repo, nil
}