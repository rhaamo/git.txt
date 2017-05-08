package repository

import (
	git "gopkg.in/libgit2/git2go.v25"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"path"
)

// Builds a repository path
func RepoPath(username string, repo_name string) string {
	return path.Join(setting.RepositoryRoot, username, repo_name) + ".git"
}

func InitRepository(username string, repo_name string) (*git.Repository, error) {
	// Join the final path of the repository
	repoPath := RepoPath(username, repo_name)

	// Create the repository, bare type
	repo, err := git.InitRepository(repoPath, true)
	if err != nil {
		return nil, err
	}
	return repo, err
}