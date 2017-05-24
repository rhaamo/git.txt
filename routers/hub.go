package routers

import "dev.sigpipe.me/dashie/git.txt/context"

func NotFound(ctx *context.Context) {
	ctx.Title(ctx.Tr("error.page_not_found"))
	ctx.Handle(404, "home.NotFound", nil)
}