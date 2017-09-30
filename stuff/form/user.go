package form

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
)

// UpdateSettingsProfile form struct
type UpdateSettingsProfile struct {
	Email string `binding:"Required;Email;MaxSize(254)"`
}

// Validate func
func (f *UpdateSettingsProfile) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}
