package form

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
)

/* New Gitxt */
type Gitxt struct {
	Description string `binding:"MaxSize(255)"`
	IsPublic    bool   `binding:"Default:1"`
	// Validation builtin into Macaron/Binding doesn't validates theses slices
	// See the router view for manual validation
	FilesFilename	[]string `binding:"Required;MaxSize(255);MinSizeSlice(1)"`
	FilesContent	[]string `binding:"Required;MaxSize(255);MinSizeSlice(1)"`
}

func (f *Gitxt) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

/* Delete Gitxt */
type GitxtDelete struct {
	Hash    string   `binding:"Required"`
	Owner    string   `binding:"Required"`
}

func (f *GitxtDelete) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}