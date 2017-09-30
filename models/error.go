// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
)

// ErrNameReserved struct
type ErrNameReserved struct {
	Name string
}

// IsErrNameReserved func
func IsErrNameReserved(err error) bool {
	_, ok := err.(ErrNameReserved)
	return ok
}

// Error func
func (err ErrNameReserved) Error() string {
	return fmt.Sprintf("name is reserved [name: %s]", err.Name)
}

// ErrNamePatternNotAllowed struct
type ErrNamePatternNotAllowed struct {
	Pattern string
}

// IsErrNamePatternNotAllowed func
func IsErrNamePatternNotAllowed(err error) bool {
	_, ok := err.(ErrNamePatternNotAllowed)
	return ok
}

// Error func
func (err ErrNamePatternNotAllowed) Error() string {
	return fmt.Sprintf("name pattern is not allowed [pattern: %s]", err.Pattern)
}

//  ____ ___
// |    |   \______ ___________
// |    |   /  ___// __ \_  __ \
// |    |  /\___ \\  ___/|  | \/
// |______//____  >\___  >__|
//              \/     \/

// ErrUserAlreadyExist struct
type ErrUserAlreadyExist struct {
	Name string
}

// IsErrUserAlreadyExist func
func IsErrUserAlreadyExist(err error) bool {
	_, ok := err.(ErrUserAlreadyExist)
	return ok
}

// Error func
func (err ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("user already exists [name: %s]", err.Name)
}

// ErrEmailAlreadyUsed struct
type ErrEmailAlreadyUsed struct {
	Email string
}

// IsErrEmailAlreadyUsed func
func IsErrEmailAlreadyUsed(err error) bool {
	_, ok := err.(ErrEmailAlreadyUsed)
	return ok
}

// Error func
func (err ErrEmailAlreadyUsed) Error() string {
	return fmt.Sprintf("e-mail has been used [email: %s]", err.Email)
}

// ErrUserOwnRepos struct
type ErrUserOwnRepos struct {
	UID int64
}

// IsErrUserOwnRepos func
func IsErrUserOwnRepos(err error) bool {
	_, ok := err.(ErrUserOwnRepos)
	return ok
}

// Error func
func (err ErrUserOwnRepos) Error() string {
	return fmt.Sprintf("user still has ownership of repositories [uid: %d]", err.UID)
}

// __________     ___.   .__  .__          ____  __.
// \______   \__ _\_ |__ |  | |__| ____   |    |/ _|____ ___.__.
//  |     ___/  |  \ __ \|  | |  |/ ___\  |      <_/ __ <   |  |
//  |    |   |  |  / \_\ \  |_|  \  \___  |    |  \  ___/\___  |
//  |____|   |____/|___  /____/__|\___  > |____|__ \___  > ____|
//                     \/             \/          \/   \/\/

// ErrKeyUnableVerify struct
type ErrKeyUnableVerify struct {
	Result string
}

// IsErrKeyUnableVerify func
func IsErrKeyUnableVerify(err error) bool {
	_, ok := err.(ErrKeyUnableVerify)
	return ok
}

// Error func
func (err ErrKeyUnableVerify) Error() string {
	return fmt.Sprintf("Unable to verify key content [result: %s]", err.Result)
}

// ErrKeyNotExist struct
type ErrKeyNotExist struct {
	ID int64
}

// IsErrKeyNotExist func
func IsErrKeyNotExist(err error) bool {
	_, ok := err.(ErrKeyNotExist)
	return ok
}

// Error func
func (err ErrKeyNotExist) Error() string {
	return fmt.Sprintf("public key does not exist [id: %d]", err.ID)
}

// ErrKeyAlreadyExist struct
type ErrKeyAlreadyExist struct {
	OwnerID int64
	Content string
}

// IsErrKeyAlreadyExist func
func IsErrKeyAlreadyExist(err error) bool {
	_, ok := err.(ErrKeyAlreadyExist)
	return ok
}

// Error func
func (err ErrKeyAlreadyExist) Error() string {
	return fmt.Sprintf("public key already exists [owner_id: %d, content: %s]", err.OwnerID, err.Content)
}

// ErrKeyNameAlreadyUsed struct
type ErrKeyNameAlreadyUsed struct {
	OwnerID int64
	Name    string
}

// IsErrKeyNameAlreadyUsed func
func IsErrKeyNameAlreadyUsed(err error) bool {
	_, ok := err.(ErrKeyNameAlreadyUsed)
	return ok
}

// Error func
func (err ErrKeyNameAlreadyUsed) Error() string {
	return fmt.Sprintf("public key already exists [owner_id: %d, name: %s]", err.OwnerID, err.Name)
}

// ErrKeyAccessDenied struct
type ErrKeyAccessDenied struct {
	UserID int64
	KeyID  int64
	Note   string
}

// IsErrKeyAccessDenied func
func IsErrKeyAccessDenied(err error) bool {
	_, ok := err.(ErrKeyAccessDenied)
	return ok
}

// Error func
func (err ErrKeyAccessDenied) Error() string {
	return fmt.Sprintf("user does not have access to the key [user_id: %d, key_id: %d, note: %s]",
		err.UserID, err.KeyID, err.Note)
}

// ErrDeployKeyNotExist struct
type ErrDeployKeyNotExist struct {
	ID     int64
	KeyID  int64
	RepoID int64
}

// IsErrDeployKeyNotExist func
func IsErrDeployKeyNotExist(err error) bool {
	_, ok := err.(ErrDeployKeyNotExist)
	return ok
}

// Error func
func (err ErrDeployKeyNotExist) Error() string {
	return fmt.Sprintf("Deploy key does not exist [id: %d, key_id: %d, repo_id: %d]", err.ID, err.KeyID, err.RepoID)
}

// ErrDeployKeyAlreadyExist struct
type ErrDeployKeyAlreadyExist struct {
	KeyID  int64
	RepoID int64
}

// IsErrDeployKeyAlreadyExist func
func IsErrDeployKeyAlreadyExist(err error) bool {
	_, ok := err.(ErrDeployKeyAlreadyExist)
	return ok
}

// Error func
func (err ErrDeployKeyAlreadyExist) Error() string {
	return fmt.Sprintf("public key already exists [key_id: %d, repo_id: %d]", err.KeyID, err.RepoID)
}

// ErrDeployKeyNameAlreadyUsed struct
type ErrDeployKeyNameAlreadyUsed struct {
	RepoID int64
	Name   string
}

// IsErrDeployKeyNameAlreadyUsed func
func IsErrDeployKeyNameAlreadyUsed(err error) bool {
	_, ok := err.(ErrDeployKeyNameAlreadyUsed)
	return ok
}

// Error func
func (err ErrDeployKeyNameAlreadyUsed) Error() string {
	return fmt.Sprintf("public key already exists [repo_id: %d, name: %s]", err.RepoID, err.Name)
}


// ErrHashAlreadyExist struct
type ErrHashAlreadyExist struct {
	Hash string
}

// IsErrHashAlreadyExist func
func IsErrHashAlreadyExist(err error) bool {
	_, ok := err.(ErrHashAlreadyExist)
	return ok
}

// Error func
func (err ErrHashAlreadyExist) Error() string {
	return fmt.Sprintf("hash already exists [hash: %s]", err.Hash)
}
