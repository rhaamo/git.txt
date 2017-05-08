package context

import (
	"gopkg.in/macaron.v1"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"github.com/go-macaron/csrf"
	"net/url"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCSRF     bool
	AnonymousCreate bool
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(ctx *Context) {
		// Redirect non-login pages from logged in user
		if options.SignOutRequired && ctx.IsLogged && ctx.Req.RequestURI != "/" {
			ctx.Redirect(setting.AppSubURL + "/")
			return
		}

		if !options.SignOutRequired && !options.DisableCSRF && ctx.Req.Method == "POST" {
			csrf.Validate(ctx.Context, ctx.csrf)
			if ctx.Written() {
				return
			}
		}

		if options.SignInRequired {
			if !ctx.IsLogged {
				ctx.SetCookie("redirect_to", url.QueryEscape(setting.AppSubURL + ctx.Req.RequestURI), 0, setting.AppSubURL)
				ctx.Redirect(setting.AppSubURL + "/user/login")
				return
			} else if !ctx.User.IsActive {
				ctx.Title("auth.activate_your_account")
				ctx.HTML(200, "user/auth/activate")
				return
			}
		}

		// Redirect to login page if auto-sign provided and not signed in
		if !options.SignOutRequired && !ctx.IsLogged && len(ctx.GetCookie(setting.CookieUserName)) > 0 {
			ctx.SetCookie("redirect_to", url.QueryEscape(setting.AppSubURL + ctx.Req.RequestURI), 0, setting.AppSubURL)
			ctx.Redirect(setting.AppSubURL + "/user/login")
			return
		}
	}
}