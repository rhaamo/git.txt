package gite

import (
	"archive/tar"
	"archive/zip"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"github.com/rakyll/magicmime"
	"gopkg.in/libgit2/git2go.v26"
	gotemplate "html/template"
	"os"
	"path/filepath"
	"strings"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// Some stuff to play with git because easy things are not always easy
// Some of theses functions comes from https://git.nurupoga.org/kr/_discorde/src/master/discorde/tree.go

// A TreeEntry represents a file of a git repository.
type TreeEntry struct {
	Path     string
	IsDir    bool
	IsParent bool
	Oid      *git.Oid
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
			*entries = append(*entries, TreeEntry{path, isDir, false, nil})
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
		*entries = append(*entries, TreeEntry{filepath.Dir(target), true, true, nil})
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

// Helper to get a Zip from repository
func getEntriesPaths(entries *[]TreeEntry, target string) git.TreeWalkCallback {
	target = strings.Trim(target, "/")
	return func(path string, entry *git.TreeEntry) int {
		path = strings.Trim(path, "/")
		path = filepath.Join(path, entry.Name)
		isDir := git.FilemodeTree == git.Filemode(entry.Filemode)
		*entries = append(*entries, TreeEntry{path, isDir, false, entry.Id})
		return 0
	}
}

func getJustRawContent(repo *git.Repository, path string) (content []byte, size int64, err error) {
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
	return blob.Contents(), blob.Size(), err
}

// WriteTarArchiveFromRepository because yes
func WriteTarArchiveFromRepository(repo *git.Repository, archivePath string) (err error) {
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	z := tar.NewWriter(archiveFile)
	defer z.Close()

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
		c, _, e := getJustRawContent(repo, entry.Path)
		if e != nil {
			return e
		}

		hdr := &tar.Header{
			Name: entry.Path,
			Mode: 0600,
			Size: int64(len(c)),
		}
		if err := z.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := z.Write(c); err != nil {
			return err
		}
	}

	return
}

// WriteZipArchiveFromRepository as said
func WriteZipArchiveFromRepository(repo *git.Repository, archivePath string) (err error) {
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	z := zip.NewWriter(archiveFile)
	defer z.Close()

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
		c, _, e := getJustRawContent(repo, entry.Path)
		if e != nil {
			return e
		}
		if _, e = f.Write(c); e != nil {
			return e
		}
	}

	return
}

// TreeFiles struct to get root tree without depth and contents
type TreeFiles struct {
	ID           string
	Path         string
	Content      string
	ContentB     []byte
	ContentH     gotemplate.HTML
	Size         int64 // bytes
	OverSize     bool
	IsBinary     bool
	OverPageSize bool
	MimeType     string
	LineNos      gotemplate.HTML
}

// Content helpers
const rawContentCheckSize = 5000

// isBinary returns true if data's format is binary.
// This function will only check the first rawContentCheckSize bytes
// so it may give false positives even if it is unlikely.
func isBinary(data []byte) bool {
	if len(data) > rawContentCheckSize {
		data = data[:rawContentCheckSize]
	}
	for _, b := range data {
		if b == byte(0x0) {
			return true
		}
	}
	return false
}
func getTreeFile(repo *git.Repository, path string, curSize int64) (treeFile TreeFiles, err error) {
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
	treeFile = TreeFiles{}

	treeFile.IsBinary = isBinary(blob.Contents())

	if len(blob.Contents()) <= 0 {
		treeFile.MimeType = "text/plain"
	} else if len(blob.Contents()) > rawContentCheckSize {
		treeFile.MimeType, err = magicmime.TypeByBuffer(blob.Contents()[:rawContentCheckSize])
	} else {
		treeFile.MimeType, err = magicmime.TypeByBuffer(blob.Contents()[:])
	}

	if err != nil {
		return
	}

	// First check if Binary
	if treeFile.IsBinary {
		treeFile.Content = ""
	} else {
		// Then the whole file size
		if blob.Size() > setting.Bloby.MaxSizeDisplay {
			treeFile.Content = ""
			treeFile.OverSize = true
		} else {
			// If we still are in the limits, check the page display size
			if curSize+blob.Size() > setting.Bloby.MaxPageDisplay {
				treeFile.OverPageSize = true
				treeFile.Content = ""
			} else {
				treeFile.OverPageSize = false
				treeFile.Content = string(blob.Contents()[:])
			}
		}
	}

	treeFile.Size = blob.Size()
	treeFile.Path = path
	treeFile.ID = blob.Id().String()

	return treeFile, err
}

// GetTreeFileNoLimit yeah
func GetTreeFileNoLimit(repo *git.Repository, path string) (treeFile TreeFiles, err error) {
	tree, err := getRepositoryTree(repo)
	if err != nil {
		return
	}
	if err = magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_ERROR); err != nil {
		return
	}
	defer magicmime.Close()

	var entry git.TreeEntry
	tree.Walk(getTreeEntryByPath(&entry, path))

	// If we don't get an entry, the file does not exist!
	if entry.Id == nil {
		return
	}

	blob, err := repo.LookupBlob(entry.Id)
	if err != nil {
		return
	}
	treeFile = TreeFiles{}

	treeFile.IsBinary = isBinary(blob.Contents())

	if len(blob.Contents()) > rawContentCheckSize {
		treeFile.MimeType, err = magicmime.TypeByBuffer(blob.Contents()[:rawContentCheckSize])
	} else {
		treeFile.MimeType, err = magicmime.TypeByBuffer(blob.Contents()[:])
	}

	if err != nil {
		return
	}

	// First check if Binary
	treeFile.OverPageSize = false
	treeFile.ContentB = blob.Contents()

	treeFile.Size = blob.Size()
	treeFile.Path = path
	treeFile.ID = blob.Id().String()

	return treeFile, err
}

// GetWalkTreeWithContent TreeFiles.Content will be nil if size is too big
func GetWalkTreeWithContent(repo *git.Repository, path string) (finalEntries []TreeFiles, err error) {
	tree, err := getRepositoryTree(repo)
	if err != nil {
		return nil, err
	}
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_ERROR); err != nil {
		return nil, err
	}
	defer magicmime.Close()

	entries := []TreeEntry{}

	tree.Walk(getEntriesPaths(&entries, path))

	var pageSize int64

	for _, entry := range entries {
		if entry.IsDir {
			continue
		}
		treeFile, err := getTreeFile(repo, entry.Path, pageSize)
		if err != nil {
			return nil, err
		}

		if !treeFile.IsBinary && !treeFile.OverSize {
			pageSize += treeFile.Size
		}

		finalEntries = append(finalEntries, treeFile)
	}
	return finalEntries, nil
}

// GetTreeFileOid returns the OID of the Tree File
func GetTreeFileOid(repo *git.Repository, path string) (oid *git.Oid, err error) {
	tree, err := getRepositoryTree(repo)
	if err != nil {
		return
	}
	var entry git.TreeEntry
	tree.Walk(getTreeEntryByPath(&entry, path))

	return entry.Id, err
}
