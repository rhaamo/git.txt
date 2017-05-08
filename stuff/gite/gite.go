package gite

import (
	"gopkg.in/libgit2/git2go.v25"
	"strings"
	"path/filepath"
	"io"
	"archive/zip"
)

// Some stuff to play with git because easy things are not always easy
// Some of theses functions comes from https://git.nurupoga.org/kr/_discorde/src/master/discorde/tree.go


// A TreeEntry represents a file of a git repository.
type TreeEntry struct {
	Path     string
	IsDir    bool
	IsParent bool
}

// appendTargetEntries returns a tree walking callback which appends
// entries located in a given target path to a slice.
func appendTargetEntries(entries *[]TreeEntry, target string) git.TreeWalkCallback {
	target = strings.Trim(target, "/")
	return func(path string, entry *git.TreeEntry) int {
		path = strings.Trim(path, "/")
		if path == target {
			isDir := git.FilemodeTree == git.Filemode(entry.Filemode)
			path = filepath.Join(path, entry.Name)
			*entries = append(*entries, TreeEntry{path, isDir, false})
		}
		return 0
	}
}

// getRepositoryTree returns the tree of a git repository starting
// at the root.
func getRepositoryTree(repo *git.Repository) (*git.Tree, error) {
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		return nil, err
	}
	return headCommit.Tree()
}

func appendParentEntry(entries *[]TreeEntry, target string) {
	if target != "/" && target != "" {
		*entries = append(*entries, TreeEntry{filepath.Dir(target), true, true})
	}
}

// getTreeEntries returns a list of TreeEntry for each file present in the
// git repository directory at the indicated path.
func getTreeEntries(repo *git.Repository, path string) (entries []TreeEntry, err error) {
	tree, err := getRepositoryTree(repo)
	if err != nil {
		return
	}
	appendParentEntry(&entries, path)
	tree.Walk(appendTargetEntries(&entries, path))
	return
}

func getTreeEntryByPath(tentry *git.TreeEntry, target string) git.TreeWalkCallback {
	target = strings.Trim(target, "/")
	return func(path string, entry *git.TreeEntry) int {
		path = strings.Trim(path, "/")
		path = filepath.Join(path, entry.Name)
		if path == target {
			*tentry = *entry
			return 0
		}
		return 0
	}
}

// Content helpers

const RAW_CONTENT_CHECK_SIZE = 5000

// isBinary returns true if data's format is binary.
// This function will only check the first RAW_CONTENT_CHECK_SIZE bytes
// so it may give false positives even if it is unlikely.
func isBinary(data []byte) bool {
	if len(data) > RAW_CONTENT_CHECK_SIZE {
		data = data[:RAW_CONTENT_CHECK_SIZE]
	}
	for _, b := range data {
		if b == byte(0x0) {
			return true
		}
	}
	return false
}
func getRawContent(repo *git.Repository, path string) (content []byte, err error) {
	tree, err := getRepositoryTree(repo)
	if err != nil {
		return
	}
	var entry git.TreeEntry
	tree.Walk(getTreeEntryByPath(&entry, path))
	blob, err := repo.LookupBlob(entry.Id)
	if err != nil {
		return
	}
	return blob.Contents(), err
}

// Helper to get a Zip from repository
func getEntriesPaths(entries *[]TreeEntry, target string) git.TreeWalkCallback {
	target = strings.Trim(target, "/")
	return func(path string, entry *git.TreeEntry) int {
		path = strings.Trim(path, "/")
		path = filepath.Join(path, entry.Name)
		isDir := git.FilemodeTree == git.Filemode(entry.Filemode)
		*entries = append(*entries, TreeEntry{path, isDir, false})
		return 0
	}
}
func WriteZipFromRepository(w io.Writer, repo *git.Repository) (err error) {
	z := zip.NewWriter(w)
	tree, err := getRepositoryTree(repo)
	if err != nil {
		return
	}
	entries := []TreeEntry{}
	tree.Walk(getEntriesPaths(&entries, "/"))
	for _, entry := range entries {
		if entry.IsDir {
			continue
		}
		f, e := z.Create(entry.Path)
		if e != nil {
			return e
		}
		c, e := getRawContent(repo, entry.Path)
		if e != nil {
			return e
		}
		if _, e = f.Write(c); e != nil {
			return e
		}
	}
	err = z.Close()
	return
}

// Get root tree without depth and contents

type TreeFiles struct {
	Path	string
	Content	string
}

func GetWalkTreeWithContent(repo *git.Repository, path string) (finalEntries []TreeFiles, err error) {
	tree, err := getRepositoryTree(repo)
	if err != nil {
		return nil, err
	}
	entries := []TreeEntry{}

	tree.Walk(getEntriesPaths(&entries, path))

	for _, entry := range entries {
		if entry.IsDir {
			continue
		}
		content, err := getRawContent(repo, entry.Path)
		if err != nil {
			return nil, err
		}
		finalEntries = append(finalEntries, TreeFiles{Path: entry.Path, Content: string(content[:])})
	}
	return finalEntries, nil
}
