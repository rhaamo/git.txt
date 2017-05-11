package routers

import "dev.sigpipe.me/dashie/git.txt/context"

func NotFound(ctx *context.Context) {
	ctx.Data["Title"] = "Page Not Found"
	ctx.Handle(404, "home.NotFound", nil)
}