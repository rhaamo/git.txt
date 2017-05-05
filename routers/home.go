package routers

import (
	"dev.sigpipe.me/dashie/git.txt/context"
)

const (
	HOME	= "home"
)

func Home(ctx *context.Context) {
	ctx.Title("git.txt")
	ctx.HTML(200, HOME)
}