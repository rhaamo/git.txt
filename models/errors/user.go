// Copyright 2017 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package errors

import "fmt"

// EmptyName struct
type EmptyName struct{}

// IsEmptyName func
func IsEmptyName(err error) bool {
	_, ok := err.(EmptyName)
	return ok
}

// Error func
func (err EmptyName) Error() string {
	return "empty name"
}

// UserNotExist struct
type UserNotExist struct {
	UserID int64
	Name   string
}

// IsUserNotExist func
func IsUserNotExist(err error) bool {
	_, ok := err.(UserNotExist)
	return ok
}

// Error func
func (err UserNotExist) Error() string {
	return fmt.Sprintf("user does not exist [user_id: %d, name: %s]", err.UserID, err.Name)
}

// UserNotKeyOwner struct
type UserNotKeyOwner struct {
	KeyID int64
}

// IsUserNotKeyOwner func
func IsUserNotKeyOwner(err error) bool {
	_, ok := err.(UserNotKeyOwner)
	return ok
}

// Error func
func (err UserNotKeyOwner) Error() string {
	return fmt.Sprintf("user is not the owner of public key [key_id: %d]", err.KeyID)
}
