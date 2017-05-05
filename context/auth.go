package context

import (
	"gopkg.in/macaron.v1"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(ctx *Context) {
		if options.SignInRequired {
		}
	}
}