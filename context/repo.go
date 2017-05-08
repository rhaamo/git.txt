package context

import (
	"gopkg.in/macaron.v1"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
)

func AssignRepository() macaron.Handler {
	return func(ctx *Context) {
		userName := ctx.Params("user")
		repoHash := ctx.Params("hash")
		repo, err := models.GetRepositoryByName(userName, repoHash)
		if err != nil {
			if errors.IsRepoNotExist(err) {
				if ctx.Query("go-get") == "1" {
					// maybe should do something here
					return
				}
				ctx.NotFound()
			} else {
				ctx.Handle(500, "GetRepositoryByName", err)
			}
			return
		}

		ctx.Gitxt = repo
		ctx.RepoOwnerUsername = userName
	}
}