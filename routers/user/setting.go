package user

import (
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/stuff/form"
	"dev.sigpipe.me/dashie/git.txt/models"
)

const (
	tmplSettingsProfile = "user/settings/profile"
)

// Settings GET
func Settings(ctx *context.Context) {
	ctx.Title("settings.title")
	ctx.PageIs("SettingsProfile")
	ctx.Data["email"] = ctx.User.Email
	ctx.Success(tmplSettingsProfile)
}

// SettingsPost POST
func SettingsPost(ctx *context.Context, f form.UpdateSettingsProfile) {
	ctx.Title("settings.title")
	ctx.PageIs("SettingsProfile")
	ctx.Data["origin_name"] = ctx.User.UserName

	if ctx.HasError() {
		ctx.Success(tmplSettingsProfile)
		return
	}

	ctx.User.Email = f.Email
	if err := models.UpdateUser(ctx.User); err != nil {
		ctx.ServerError("UpdateUser", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("settings.update_profile_success"))
	ctx.SubURLRedirect("/user/settings")
}