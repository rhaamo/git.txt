package errors

import "fmt"

// RepoNotExist struct
type RepoNotExist struct {
	ID     int64
	UserID int64
	Name   string
}

// IsRepoNotExist func
func IsRepoNotExist(err error) bool {
	_, ok := err.(RepoNotExist)
	return ok
}

// Error func
func (err RepoNotExist) Error() string {
	return fmt.Sprintf("repository does not exist [id: %d, user_id: %d, name: %s]", err.ID, err.UserID, err.Name)
}
