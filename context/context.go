package context

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/go-macaron/cache"
	log "gopkg.in/clog.v1"
	"net/http"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"fmt"
	"io"
	"time"
	"strings"
	"github.com/go-macaron/i18n"
	"html/template"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/stuff/form"
	"dev.sigpipe.me/dashie/git.txt/stuff/auth"
)

// Context represents context of a request.
type Context struct {
	*macaron.Context
	Cache   cache.Cache
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	User	*models.User
	Gitxt   *models.Gitxt

	RepoOwnerUsername  string

	IsLogged	bool
	IsBasicAuth	bool
}

// Title sets "Title" field in template data.
func (c *Context) Title(locale string) {
	c.Data["Title"] = c.Tr(locale) + " - " + setting.AppName
}

// PageIs sets "PageIsxxx" field in template data.
func (c *Context) PageIs(name string) {
	c.Data["PageIs"+name] = true
}

// HTML responses template with given status.
func (ctx *Context) HTML(status int, name string) {
	log.Trace("Template: %s", name)
	ctx.Context.HTML(status, name)
}

// Success responses template with status http.StatusOK.
func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

// JSONSuccess responses JSON with status http.StatusOK.
func (c *Context) JSONSuccess(data interface{}) {
	c.JSON(http.StatusOK, data)
}

// HasError returns true if error occurs in form validation.
func (ctx *Context) HasError() bool {
	hasErr, ok := ctx.Data["HasError"]
	if !ok {
		return false
	}
	ctx.Flash.ErrorMsg = ctx.Data["ErrorMsg"].(string)
	ctx.Data["Flash"] = ctx.Flash
	return hasErr.(bool)
}

// RenderWithErr used for page has form validation but need to prompt error to users.
func (ctx *Context) RenderWithErr(msg, tpl string, f interface{}) {
	if f != nil {
		form.Assign(f, ctx.Data)
	}
	ctx.Flash.ErrorMsg = msg
	ctx.Data["Flash"] = ctx.Flash
	ctx.HTML(http.StatusOK, tpl)
}

// Handle handles and logs error by given status.
func (ctx *Context) Handle(status int, title string, err error) {
	switch status {
	case http.StatusNotFound:
		ctx.Data["Title"] = "Page Not Found"
	case http.StatusInternalServerError:
		ctx.Data["Title"] = "Internal Server Error"
		log.Error(2, "%s: %v", title, err)
	}
	ctx.HTML(status, fmt.Sprintf("status/%d", status))
}

// NotFound renders the 404 page.
func (ctx *Context) NotFound() {
	ctx.Handle(http.StatusNotFound, "", nil)
}

// ServerError renders the 500 page.
func (c *Context) ServerError(title string, err error) {
	c.Handle(http.StatusInternalServerError, title, err)
}

// SubURLRedirect responses redirection wtih given location and status.
// It prepends setting.AppSubURL to the location string.
func (c *Context) SubURLRedirect(location string, status ...int) {
	c.Redirect(setting.AppSubURL + location)
}

// NotFoundOrServerError use error check function to determine if the error
// is about not found. It responses with 404 status code for not found error,
// or error context description for logging purpose of 500 server error.
func (c *Context) NotFoundOrServerError(title string, errck func(error) bool, err error) {
	if errck(err) {
		c.NotFound()
		return
	}
	c.ServerError(title, err)
}

func (ctx *Context) ServeContent(name string, r io.ReadSeeker, params ...interface{}) {
	modtime := time.Now()
	for _, p := range params {
		switch v := p.(type) {
		case time.Time:
			modtime = v
		}
	}
	ctx.Resp.Header().Set("Content-Description", "File Transfer")
	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+name)
	ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Resp.Header().Set("Expires", "0")
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.Header().Set("Pragma", "public")
	http.ServeContent(ctx.Resp, ctx.Req.Request, name, modtime, r)
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(c *macaron.Context, l i18n.Locale, cache cache.Cache, sess session.Store, f *session.Flash, x csrf.CSRF) {
		ctx := &Context{
			Context: c,
			Cache:   cache,
			csrf:    x,
			Flash:   f,
			Session: sess,
		}

		if len(setting.HTTP.AccessControlAllowOrigin) > 0 {
			ctx.Header().Set("Access-Control-Allow-Origin", setting.HTTP.AccessControlAllowOrigin)
			ctx.Header().Set("'Access-Control-Allow-Credentials' ", "true")
			ctx.Header().Set("Access-Control-Max-Age", "3600")
			ctx.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		}

		// Compute current URL for real-time change language.
		ctx.Data["Link"] = setting.AppSubURL + strings.TrimSuffix(ctx.Req.URL.Path, "/")

		ctx.Data["PageStartTime"] = time.Now()

		// Get user from session if logined.
		ctx.User, ctx.IsBasicAuth = auth.SignedInUser(ctx.Context, ctx.Session)

		if ctx.User != nil {
			ctx.IsLogged = true
			ctx.Data["IsLogged"] = ctx.IsLogged
			ctx.Data["UserIsAdmin"] = ctx.User.IsAdmin
			ctx.Data["LoggedUser"] = ctx.User
			ctx.Data["LoggedUserID"] = ctx.User.ID
			ctx.Data["LoggedUserName"] = ctx.User.UserName
		} else {
			ctx.Data["LoggedUserID"] = 0
			ctx.Data["LoggedUserName"] = ""
		}

		ctx.Data["CSRFToken"] = x.GetToken()
		ctx.Data["CSRFTokenHTML"] = template.HTML(`<input type="hidden" name="_csrf" value="` + x.GetToken() + `">`)
		log.Trace("Session ID: %s", sess.ID())
		log.Trace("CSRF Token: %v", ctx.Data["CSRFToken"])

		c.Map(ctx)
	}
}