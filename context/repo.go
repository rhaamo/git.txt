package context

import (
	"gopkg.in/macaron.v1"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
)

type Gitxt struct {
	User	*models.User
	Gitxt	*models.Gitxt
	Owner	bool
	UserName	string
}

func AssignRepository() macaron.Handler {
	return func(ctx *Context) {
		userName := ctx.Params("user")
		repoHash := ctx.Params("hash")
		repo, err := models.GetRepositoryByName(userName, repoHash)
		if err != nil {
			if errors.IsRepoNotExist(err) {
				ctx.NotFound()
			} else {
				ctx.Handle(500, "GetRepositoryByName", err)
			}
			return
		}

		ctx.Gitxt.Gitxt = repo
		ctx.Gitxt.UserName = userName
		if ctx.Gitxt.Gitxt.Anonymous {
			ctx.Gitxt.Owner = false
		} else {
			ctx.Gitxt.Owner = ctx.Gitxt.User.ID == ctx.Gitxt.Gitxt.UserID
		}
	}
}
