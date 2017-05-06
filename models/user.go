package models

import (
	"time"
	"strings"
	"unicode/utf8"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"fmt"
	"crypto/subtle"
	"dev.sigpipe.me/dashie/git.txt/stuff/tool"
	"os"
	"path/filepath"
	"dev.sigpipe.me/dashie/git.txt/setting"
)

type User struct {
	ID		int64	`xorm:"pk autoincr"`
	UserName	string	`xorm:"UNIQUE NOT NULL"`
	LowerName	string	`xorm:"UNIQUE NOT NULL"`
	Email		string	`xorm:"NOT NULL"`

	Password	string	`xorm:"NOT NULL"`
	Rands		string	`xorm:"VARCHAR(10)"`
	Salt		string	`xorm:"VARCHAR(10)"`

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

func countUsers(e Engine) int64 {
	count, _ := x.Count(new(User))
	return count
}

// CountUsers returns number of users.
func CountUsers() int64 {
	return countUsers(x)
}

func getUserByID(e Engine, id int64) (*User, error) {
	u := new(User)
	has, err := e.Id(id).Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.UserNotExist{id, ""}
	}
	return u, nil
}

// GetUserByID returns the user object by given ID if exists.
func GetUserByID(id int64) (*User, error) {
	return getUserByID(x, id)
}

// IsUserExist checks if given user name exist,
// the user name should be noncased unique.
// If uid is presented, then check will rule out that one,
// it is used when update a user name in settings page.
func IsUserExist(uid int64, name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return x.Where("id != ?", uid).Get(&User{LowerName: strings.ToLower(name)})
}

var (
	reservedUsernames    = []string{"assets", "css", "img", "js", "less", "plugins", "debug", "raw", "install", "api", "avatar", "user", "org", "help", "stars", "issues", "pulls", "commits", "repo", "template", "admin", "new", ".", ".."}
	reservedUserPatterns = []string{"*.keys"}
)

// isUsableName checks if name is reserved or pattern of name is not allowed
// based on given reserved names and patterns.
// Names are exact match, patterns can be prefix or suffix match with placeholder '*'.
func isUsableName(names, patterns []string, name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return errors.EmptyName{}
	}

	for i := range names {
		if name == names[i] {
			return ErrNameReserved{name}
		}
	}

	for _, pat := range patterns {
		if pat[0] == '*' && strings.HasSuffix(name, pat[1:]) ||
			(pat[len(pat)-1] == '*' && strings.HasPrefix(name, pat[:len(pat)-1])) {
			return ErrNamePatternNotAllowed{pat}
		}
	}

	return nil
}

func IsUsableUsername(name string) error {
	return isUsableName(reservedUsernames, reservedUserPatterns, name)
}

// EncodePasswd encodes password to safe format.
func (u *User) EncodePasswd() {
	newPasswd := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(passwd string) bool {
	newUser := &User{Password: passwd, Salt: u.Salt}
	newUser.EncodePasswd()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}

// GetUserSalt returns a ramdom user salt token.
func GetUserSalt() (string, error) {
	return tool.RandomString(10)
}

// UserPath returns the path absolute path of user repositories.
func UserPath(userName string) string {
	return filepath.Join(setting.RepositoryRoot, strings.ToLower(userName))
}

// Create a new user and do some validation
func CreateUser(u *User) (err error) {
	if err = IsUsableUsername(u.UserName); err != nil {
		return err
	}

	isExist, err := IsUserExist(0, u.UserName)
	if err != nil {
		return err
	} else if isExist {
		return ErrUserAlreadyExist{u.UserName}
	}

	u.Email = strings.ToLower(u.Email)
	u.LowerName = strings.ToLower(u.UserName)

	if u.Rands, err = GetUserSalt(); err != nil {
		return err
	}
	if u.Salt, err = GetUserSalt(); err != nil {
		return err
	}
	u.EncodePasswd()

	sess := x.NewSession()
	defer sessionRelease(sess)
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Insert(u); err != nil {
		return err
	} else if err = os.MkdirAll(UserPath(u.UserName), os.ModePerm); err != nil {
		return err
	}

	return sess.Commit()
}

// Update an user
func updateUser(e Engine, u *User) error {
	u.LowerName = strings.ToLower(u.UserName)
	u.Email = strings.ToLower(u.Email)
	_, err := e.Id(u.ID).AllCols().Update(u)
	return err
}

func UpdateUser(u *User) error {
	return updateUser(x, u)
}

// Login validates user name and password.
func UserLogin(username, password string) (*User, error) {
	var user *User
	if strings.Contains(username, "@") {
		user = &User{Email: strings.ToLower(username)}
	} else {
		user = &User{LowerName: strings.ToLower(username)}
	}

	hasUser, err := x.Get(user)
	if err != nil {
		return nil, err
	}

	if hasUser {
		if user.ValidatePassword(password) {
			return user, nil
		}

		return nil, errors.UserNotExist{user.ID, user.UserName}
	}

	return nil, errors.UserNotExist{user.ID, user.UserName}
}