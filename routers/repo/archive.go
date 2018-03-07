package repo

import (
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/stuff/repository"
	log "gopkg.in/clog.v1"

	"dev.sigpipe.me/dashie/git.txt/stuff/gite"
	"dev.sigpipe.me/dashie/git.txt/stuff/tool"
	"github.com/Unknwon/com"
	"gopkg.in/libgit2/git2go.v25"
	"os"
	"path"
	"strings"
)

// DownloadArchive of repository
func DownloadArchive(ctx *context.Context) {
	var (
		uri         = ctx.Params("*")
		refName     string
		ext         string
		archivePath string
		archiveType string
	)

	pathRepo := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Gitxt.Hash)

	switch {
	case strings.HasSuffix(uri, ".zip"):
		ext = ".zip"
		archivePath = path.Join(pathRepo, "archives/zip")
		archiveType = "zip"
	case strings.HasSuffix(uri, ".tar.gz"):
		ext = ".tar.gz"
		archivePath = path.Join(pathRepo, "archives/targz")
		archiveType = "tar.gz"
	default:
		log.Trace("Unknown format: %s", uri)
		ctx.Error(404)
		return
	}

	refName = strings.TrimSuffix(uri, ext)
	if refName != "master" {
		ctx.Handle(500, "Ref other than master not supported", nil)
		return
	}
	// TODO: add support for != master ref

	if !com.IsDir(archivePath) {
		if err := os.MkdirAll(archivePath, os.ModePerm); err != nil {
			ctx.Handle(500, "Download -> os.MkdirAll(archivePath)", err)
			return
		}
	}

	// Get latest commit
	repo, err := git.OpenRepository(pathRepo)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty()
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	// Get repo head
	repoHead, err := repo.Head()
	if err != nil {
		log.Warn("git_error_get_head: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_head")
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	// Get latest commit
	headCommit, err := repo.LookupCommit(repoHead.Target())
	if err != nil {
		log.Warn("git_error_get_head_commit: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_head_commit")
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	archivePath = path.Join(archivePath, tool.ShortSHA1(headCommit.Object.Id().String())+ext)

	log.Trace("Going to create archive: %s", archivePath)

	if !com.IsFile(archivePath) {
		if archiveType == "zip" {
			if err := gite.WriteZipArchiveFromRepository(repo, archivePath); err != nil {
				os.RemoveAll(archivePath)
				ctx.Handle(500, "Download -> WriteArchiveFromRepository "+archivePath, err)
				return
			}
		} else if archiveType == "tar.gz" {
			if err := gite.WriteTarArchiveFromRepository(repo, archivePath); err != nil {
				os.RemoveAll(archivePath)
				ctx.Handle(500, "Download -> WriteArchiveFromRepository "+archivePath, err)
				return
			}
		} else {
			ctx.Handle(500, "Cannot generate archive type: %s", nil)
		}
	}

	ctx.ServeFile(archivePath, ctx.Gitxt.Gitxt.Hash+"-"+refName+ext)
}
