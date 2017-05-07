package gitxt

import (
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/stuff/form"
	"dev.sigpipe.me/dashie/git.txt/stuff/sanitize"
	// "dev.sigpipe.me/dashie/git.txt/models"
	"fmt"
	"strings"
)

const (
	NEW = "gitxt/new"
)

func New(ctx *context.Context) {
	ctx.Title("gitxt_new.title")
	ctx.PageIs("GitxtNew")

	// One initial file
	ctx.Data["FilesContent"] = []string{""}
	ctx.Data["FilesFilename"] = []string{""}

	ctx.Success(NEW)
}

func NewPost(ctx *context.Context, f form.Gitxt) {
	ctx.Title("gitxt_new.title")
	ctx.PageIs("GitxtNewPost")

	if ctx.HasError() {
		ctx.Success(NEW)
		return
	}

	for i := range f.FilesFilename {
		// For each filename sanitize it
		f.FilesFilename[i] = sanitize.SanitizeFilename(f.FilesFilename[i])
		if len(f.FilesFilename[i]) == 0  || f.FilesFilename[i] == "." {
			// If length is zero, use default filename
			f.FilesFilename[i] = fmt.Sprintf("gitxt%d.txt", i)
		}
	}

	for i := range f.FilesContent {
		if len(strings.TrimSpace(f.FilesContent[i])) <= 0 {
			ctx.Data[fmt.Sprintf("Err_FilesContent_%d", i)] = ctx.Tr("gitxt_new.error_files_content")
			ctx.Data["HasError"] = true
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.files_content_cannot_empty")
		}
	}

	// Since the validation of slices doesn't works in bindings, manually update context
	ctx.Data["FilesFilename"] = f.FilesFilename
	ctx.Data["FilesContent"] = f.FilesContent

	// We got an error in the manual validation step, render with error
	if ctx.HasError() {
		//ctx.RenderWithErr(ctx.Tr("gitxt_new.error_plz_correct"), NEW, &f)
		ctx.Success(NEW)
		return
	}

	// Ok we are good

	// 1. Init a bare repository

	// 2. Create the files

	// 3. Commit

	// 4. Insert info in database

	// 5. Return render to gitxt view page

	// u := &models.User{
	// 	UserName:	f.UserName,
	// 	Email:		f.Email,
	// 	Password:	f.Password,
	// 	IsActive:	true, // FIXME: implement user activation by email
	// }
	// if err := models.CreateUser(u); err != nil {
	// 	switch {
	// 	case models.IsErrUserAlreadyExist(err):
	// 		ctx.Data["Err_UserName"] = true
	// 		ctx.RenderWithErr(ctx.Tr("form.username_been_taken"), REGISTER, &f)
	// 	case models.IsErrNameReserved(err):
	// 		ctx.Data["Err_UserName"] = true
	// 		ctx.RenderWithErr(ctx.Tr("form.username_reserved"), REGISTER, &f)
	// 	case models.IsErrNamePatternNotAllowed(err):
	// 		ctx.Data["Err_UserName"] = true
	// 		ctx.RenderWithErr(ctx.Tr("form.username_pattern_not_allowed"), REGISTER, &f)
	// 	default:
	// 		ctx.Handle(500, "CreateUser", err)
	// 	}
	// 	return
	// }
	// log.Trace("Account created: %s", u.UserName)
//
	// // Auto set Admin if first user
	// if models.CountUsers() == 1 {
	// 	u.IsAdmin = true
	// 	u.IsActive = true // bypass email activation
	// 	if err := models.UpdateUser(u); err != nil {
	// 		ctx.Handle(500, "UpdateUser", err)
	// 		return
	// 	}
	// }
	//
	// ctx.Redirect(setting.AppSubURL + "/user/login")
}