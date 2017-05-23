package models

import (
	"time"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
	"github.com/go-xorm/xorm"
	"fmt"
	"os"
	"dev.sigpipe.me/dashie/git.txt/stuff/repository"
	"dev.sigpipe.me/dashie/git.txt/stuff/sync"
	"dev.sigpipe.me/dashie/git.txt/setting"
	log "gopkg.in/clog.v1"
	"path/filepath"
	"github.com/Unknwon/com"
	"strings"
	"os/exec"
)

type Gitxt struct {
	ID          int64        `xorm:"pk autoincr"`
	Hash        string        `xorm:"UNIQUE NOT NULL"`
	UserID      int64        `xorm:"INDEX"`
	Anonymous   bool
	Description string        `xorm:"TEXT"`

	// Choosen expiry in hours
	Expiry		int64
	// Calculated expiry unix timestamp from the time of creation/update
	ExpiryUnix	int64	`xorm:"INDEX"`

	// Permissions
	IsPrivate bool        `xorm:"DEFAULT 0"`

	Created     time.Time        `xorm:"-"`
	CreatedUnix int64
	Updated     time.Time        `xorm:"-"`
	UpdatedUnix int64

	// Relations
	// 	UserID
}

func (gitxt *Gitxt) BeforeInsert() {
	gitxt.CreatedUnix = time.Now().Unix()
	gitxt.UpdatedUnix = gitxt.CreatedUnix
}

func (gitxt *Gitxt) BeforeUpdate() {
	gitxt.UpdatedUnix = time.Now().Unix()
}

type GitxtWithUser struct {
	User        `xorm:"extends"`
	Gitxt        `xorm:"extends"`
}

// Preventing duplicate running tasks
var taskStatusTable = sync.NewStatusTable()

const (
	_CLEAN_OLD_ARCHIVES = "clean_old_archives"
)

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
	if err != nil && user != "anonymous" {
		return nil, err
	}

	repo := &Gitxt{
		Hash: name,
	}

	if user == "anonymous" {
		repo.UserID = 0
	} else {
		repo.UserID = u.ID
	}

	has, err := x.Get(repo)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.RepoNotExist{0, repo.UserID, name}
	}
	return repo, nil
}

type GitxtOptions struct {
	UserID      int64
	WithPrivate bool
	GetAll      bool
	Page        int
	PageSize    int
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

// Archive deletion
func DeleteOldRepositoryArchives() {
	if taskStatusTable.IsRunning(_CLEAN_OLD_ARCHIVES) {
		return
	}
	taskStatusTable.Start(_CLEAN_OLD_ARCHIVES)
	defer taskStatusTable.Stop(_CLEAN_OLD_ARCHIVES)

	log.Trace("Doing: DeleteOldRepositoryArchives")

	formats := []string{"zip", "targz"}
	oldestTime := time.Now().Add(-setting.Cron.RepoArchiveCleanup.OlderThan)

	if err := x.Where("gitxt.id > 0").Table(&Gitxt{}).Join("LEFT", "user", "gitxt.user_id = user.id").Iterate(new(GitxtWithUser),
		func(idx int, bean interface{}) error {
			repo := bean.(*GitxtWithUser)
			basePath := filepath.Join(repository.RepoPath(repo.User.UserName, repo.Gitxt.Hash), "archives")
			for _, format := range formats {
				dirPath := filepath.Join(basePath, format)
				if !com.IsDir(dirPath) {
					continue
				}

				dir, err := os.Open(dirPath)
				if err != nil {
					log.Error(3, "Fail to open directory '%s': %v", dirPath, err)
					continue
				}

				fis, err := dir.Readdir(0)
				dir.Close()
				if err != nil {
					log.Error(3, "Fail to read directory '%s': %v", dirPath, err)
					continue
				}

				for _, fi := range fis {
					if fi.IsDir() || fi.ModTime().After(oldestTime) {
						continue
					}

					archivePath := filepath.Join(dirPath, fi.Name())
					if err = os.Remove(archivePath); err != nil {
						desc := fmt.Sprintf("Fail to health delete archive '%s': %v", archivePath, err)
						log.Warn(desc)
					}
				}
			}

			return nil
		}); err != nil {
		log.Error(2, "DeleteOldRepositoryArchives: %v", err)
	}
}

func removeRepository(path string) {
	var err error
	// LEGACY [Go 1.7]: workaround for Go not being able to remove read-only files/folders: https://github.com/golang/go/issues/9606
	// this bug should be fixed on Go 1.7, so the workaround should be removed when Gogs don't support Go 1.6 anymore:
	// https://github.com/golang/go/commit/2ffb3e5d905b5622204d199128dec06cefd57790
	// Note: Windows complains when delete target does not exist, therefore we can skip deletion in such cases.
	if setting.IsWindows && com.IsExist(path) {
		// converting "/" to "\" in path on Windows
		path = strings.Replace(path, "/", "\\", -1)
		err = exec.Command("cmd", "/C", "rmdir", "/S", "/Q", path).Run()
	} else {
		err = os.RemoveAll(path)
	}

	if err != nil {
		desc := fmt.Sprintf("%s: %v", path, err)
		log.Warn(desc)
	}
}

// Delete repository :'(
func DeleteRepository(ownerID int64, repoID int64) error {
	repo := &Gitxt{ID: repoID, UserID: ownerID}
	has, err := x.Get(repo)
	if err != nil {
		return err
	} else if !has {
		return errors.RepoNotExist{repoID, ownerID, ""}
	}

	repoUser, err := GetUserByID(ownerID)
	if err != nil {
		return fmt.Errorf("GetUserByID: %v", err)
	}

	sess := x.NewSession()
	defer sessionRelease(sess)
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Delete(&Gitxt{ID: repoID}); err != nil {
		return fmt.Errorf("sess.Delete: %v", err)
	}

	if err = sess.Commit(); err != nil {
		return fmt.Errorf("Commit: %v", err)
	}

	pathRepo := repository.RepoPath(repoUser.UserName, repo.Hash)
	removeRepository(pathRepo)

	return nil
}
