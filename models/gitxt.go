package models

import (
	"dev.sigpipe.me/dashie/git.txt/models/errors"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"dev.sigpipe.me/dashie/git.txt/stuff/repository"
	"dev.sigpipe.me/dashie/git.txt/stuff/sync"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/go-xorm/xorm"
	log "gopkg.in/clog.v1"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Gitxt struct
type Gitxt struct {
	ID          int64  `xorm:"pk autoincr"`
	Hash        string `xorm:"UNIQUE NOT NULL"`
	UserID      int64  `xorm:"INDEX"`
	Anonymous   bool
	Description string `xorm:"TEXT"`

	// Chosen expiry in hours
	ExpiryHours int64 `xorm:"INDEX"`
	// Calculated expiry unix timestamp from the time of creation/update
	ExpiryUnix int64
	Expiry     time.Time `xorm:"-"`

	// Permissions
	IsPrivate bool `xorm:"DEFAULT 0"`

	Created     time.Time `xorm:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-"`
	UpdatedUnix int64

	// Relations
	// 	UserID
}

// BeforeInsert hooks
func (gitxt *Gitxt) BeforeInsert() {
	gitxt.CreatedUnix = time.Now().Unix()
	gitxt.UpdatedUnix = gitxt.CreatedUnix
}

// BeforeUpdate hooks
func (gitxt *Gitxt) BeforeUpdate() {
	gitxt.UpdatedUnix = time.Now().Unix()
}

// AfterSet hooks
func (gitxt *Gitxt) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		gitxt.Created = time.Unix(gitxt.CreatedUnix, 0).Local()
	case "updated_unix":
		gitxt.Updated = time.Unix(gitxt.UpdatedUnix, 0).Local()
	case "expiry_unix":
		gitxt.Expiry = time.Unix(gitxt.ExpiryUnix, 0).Local()
	}
}

// GitxtWithUser struct
type GitxtWithUser struct {
	User  `xorm:"extends"`
	Gitxt `xorm:"extends"`
}

// Preventing duplicate running tasks
var taskStatusTable = sync.NewStatusTable()

const (
	cleanOldArchives          = "clean_old_archives"
	deleteExpiredRepositories = "delete_expired_repositories"
)

// IsHashUsed checks if given hash exist,
func IsHashUsed(uid int64, hash string) (bool, error) {
	if len(hash) == 0 {
		return false, nil
	}
	return x.Get(&Gitxt{Hash: hash})
}

// CreateGitxt Create a new gitxt
func CreateGitxt(g *Gitxt) (err error) {
	isExist, err := IsHashUsed(0, g.Hash)
	if err != nil {
		return err
	} else if isExist {
		return ErrHashAlreadyExist{g.Hash}
	}

	sess := x.NewSession()
	defer sess.Close()
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
		return nil, errors.RepoNotExist{ID: 0, UserID: repo.UserID, Name: name}
	}
	return repo, nil
}

// GitxtOptions struct
type GitxtOptions struct {
	UserID      int64
	WithPrivate bool
	GetAll      bool
	Page        int
	PageSize    int
}

// GetGitxts Get gitxts
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

// UpdateGitxt with infos
func UpdateGitxt(u *Gitxt) error {
	return updateGitxt(x, u)
}

// DeleteExpiredRepositories Delete expired
func DeleteExpiredRepositories() {
	if taskStatusTable.IsRunning(deleteExpiredRepositories) {
		return
	}
	taskStatusTable.Start(deleteExpiredRepositories)
	defer taskStatusTable.Stop(deleteExpiredRepositories)

	log.Trace("Doing: DeleteExpiredRepositories")

	type GitxtExpired struct {
		userID int64
		repoID int64
		hash   string
	}
	expired := []GitxtExpired{}

	if err := x.Where("expiry_hours > 0").And("expiry_unix <= ?", time.Now().Unix()).Iterate(new(Gitxt),
		func(idx int, bean interface{}) error {
			repo := bean.(*Gitxt)

			log.Trace("Deleting expired repository: %d/%s", repo.UserID, repo.Hash)

			expired = append(expired, GitxtExpired{repo.UserID, repo.ID, repo.Hash})

			return nil
		}); err != nil {
		log.Error(2, "DeleteExpiredRepositories: %v", err)
	}

	for _, tc := range expired {
		err := DeleteRepository(tc.userID, tc.repoID)
		if err != nil {
			log.Warn("Error removing repository %d/%d: %v", tc.userID, tc.repoID, err)
		} else {
			counter, _ := GetCounterGitxts()
			if counter.Count <= 0 {
				// This should not happens but well, anyway
				UpdateCounterGitxts(counter.Count)
			} else {
				UpdateCounterGitxts(counter.Count - 1)
			}
			log.Trace("Deleted repository %s", tc.hash)
		}

	}
}

// DeleteOldRepositoryArchives Archive deletion
func DeleteOldRepositoryArchives() {
	if taskStatusTable.IsRunning(cleanOldArchives) {
		return
	}
	taskStatusTable.Start(cleanOldArchives)
	defer taskStatusTable.Stop(cleanOldArchives)

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

// DeleteRepository Delete repository :'(
func DeleteRepository(ownerID int64, repoID int64) error {
	repo := &Gitxt{ID: repoID, UserID: ownerID}
	has, err := x.Get(repo)
	if err != nil {
		return err
	} else if !has {
		return errors.RepoNotExist{ID: repoID, UserID: ownerID, Name: ""}
	}

	// By defaults use anonymoyus
	username := "anonymous"
	if ownerID > 0 {
		// If it's a non-anonymous used, fetch it and set username to the user username
		repoUser, err := GetUserByID(ownerID)
		if err != nil {
			return fmt.Errorf("GetUserByID: %v", err)
		}
		username = repoUser.UserName
	}

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Delete(&Gitxt{ID: repoID}); err != nil {
		return fmt.Errorf("sess.Delete: %v", err)
	}

	if err = sess.Commit(); err != nil {
		return fmt.Errorf("Commit: %v", err)
	}

	pathRepo := repository.RepoPath(username, repo.Hash)
	removeRepository(pathRepo)

	return nil
}
