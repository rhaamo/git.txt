package repo

import (
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/stuff/repository"
	log "gopkg.in/clog.v1"

	"gopkg.in/libgit2/git2go.v25"
	"dev.sigpipe.me/dashie/git.txt/stuff/gite"
	"bytes"
)

func DownloadArchive(ctx *context.Context) {
	pathRepo := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Gitxt.Hash)

	repo, err := git.OpenRepository(pathRepo)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: could not open repository.")
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty();
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: repository is empty or corrupted")
		ctx.Handle(500, "GitxtView", err)
		return
	}

	var zip bytes.Buffer
	err = gite.WriteZipFromRepository(&zip, repo)
	if err != nil {
		ctx.Handle(500, "WriteZipFromRepository", err)
		return
	}

	filename := ctx.Gitxt.Gitxt.Hash + "-master.zip"
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+filename)
	ctx.Resp.Header().Set("Content-Type", "application/zip")

	zip.WriteTo(ctx.Resp)
	return
}