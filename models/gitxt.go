package models

import (
	"time"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
	"github.com/go-xorm/xorm"
	"fmt"
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

type GitxtWithUser struct {
	User	`xorm:"extends"`
	Gitxt	`xorm:"extends"`
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

type GitxtOptions struct {
	UserID	int64
	WithPrivate	bool
	GetAll		bool
	Page		int
	PageSize	int
}

// Get gitxts
func GetGitxts(opts *GitxtOptions) (gitxts []*GitxtWithUser, _ int64, _ error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	gitxts = make([]*GitxtWithUser, 0, opts.PageSize)

	sess := x.Where("is_private=?", false)

	if opts.WithPrivate && !opts.GetAll {
		sess.Or("is_private=?", true)
	}

	if !opts.GetAll {
		sess.And("user_id=?", opts.UserID)
	}

	sess.Desc("gitxt.updated_unix")

	var countSess xorm.Session
	countSess = *sess
	count, err := countSess.Count(new(Gitxt))
	if err != nil {
		return nil, 0, fmt.Errorf("Count: %v", err)
	}

	sess.Table(&Gitxt{}).Join("LEFT", "user", "gitxt.user_id = user.id")

	return gitxts, count, sess.Limit(opts.PageSize, (opts.Page-1)*opts.PageSize).Find(&gitxts)
}

// Update a Gitxt
func updateGitxt(e Engine, u *Gitxt) error {
	_, err := e.Id(u.ID).AllCols().Update(u)
	return err
}

func UpdateGitxt(u *Gitxt) error {
	return updateGitxt(x, u)
}
