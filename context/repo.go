package context

import (
	"gopkg.in/macaron.v1"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
	"strings"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"time"
	log "gopkg.in/clog.v1"
)

// Gitxt struct
type Gitxt struct {
	User	*models.User
	Gitxt	*models.Gitxt
	Owner	bool		// Gitxt.Owner: "toto/repos" does "repos" belongs to "toto" ?
	UserName	string
}

// AssignRepository to context
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

// GitUACheck User Agent check
func GitUACheck() macaron.Handler {
	return func(ctx *Context) {
		if strings.HasPrefix(strings.Join(ctx.Req.Header["User-Agent"], ""), "git/") {
			ctx.Redirect(setting.AppSubURL + strings.Replace(ctx.Req.URL.String(), ctx.Params("hash"), ctx.Params("hash")+".git", 1))
		} else {
			ctx.Redirect(setting.AppSubURL + "/" + ctx.Params("user") + "/" + ctx.Params("hash"))
		}
	}
}

// CheckRepoExpiry and handle
func CheckRepoExpiry() macaron.Handler {
	return func(ctx *Context) {
		if ctx.Gitxt.Gitxt.ExpiryHours > 0 && ctx.Gitxt.Gitxt.ExpiryUnix <= time.Now().Unix() {
			log.Trace("Repository %s/%s expired", ctx.Gitxt.UserName, ctx.Gitxt.Gitxt.Hash)
			ctx.Handle(404, "RepositoryExpired", nil)
			//ctx.Redirect(setting.AppSubURL + "/", 404)
			return
		}
	}
}