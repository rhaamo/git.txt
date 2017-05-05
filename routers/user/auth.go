package user

import (
	"dev.sigpipe.me/dashie/git.txt/context"
)

const (
	LOGIN	= "user/auth/login"
	REGISTER = "user/auth/register"
	ACTIVATE = "user/auth/activate"
	FORGOT_PASSWORD = "user/auth/forgot_password"
	RESET_PASSWORD = "user/auth/reset_password"
)

func Login(ctx *context.Context) {
	ctx.Title("login.title")

	// TODO: auto login remember_me

	ctx.HTML(200, LOGIN)
}

func Register(ctx *context.Context) {
	ctx.Title("register.title")
	ctx.HTML(200, REGISTER)
}