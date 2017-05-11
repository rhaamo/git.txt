package context

import (
	"gopkg.in/macaron.v1"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
)

func AssignUser() macaron.Handler {
	return func(ctx *Context) {
		userName := ctx.Params("user")

		// Anonymous user doesn't really exists, that's nil
		if userName == "anonymous" {
			ctx.Gitxt.User = nil
			ctx.RepoOwnerUsername = "anonymous"
			return
		}

		user, err := models.GetUserByName(userName)
		if err != nil {
			if errors.IsUserNotExist(err) {
				ctx.NotFound()
			} else {
				ctx.Handle(500, "GetUserByName", err)
			}
			return
		}

		ctx.Gitxt.User = user
		ctx.RepoOwnerUsername = userName
	}
}