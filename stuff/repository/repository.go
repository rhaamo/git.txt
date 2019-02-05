package repository

import (
	"dev.sigpipe.me/dashie/git.txt/setting"
	git "gopkg.in/libgit2/git2go.v26"
	"path"
)

// RepoPath will builds a repository path
func RepoPath(username string, repoName string) string {
	return path.Join(setting.RepositoryRoot, username, repoName) + ".git"
}

// InitRepository will init a repository with given informations
func InitRepository(username string, repoName string) (*git.Repository, error) {
	// Join the final path of the repository
	repoPath := RepoPath(username, repoName)

	// Create the repository, bare type
	repo, err := git.InitRepository(repoPath, true)
	if err != nil {
		return nil, err
	}
	return repo, err
}
