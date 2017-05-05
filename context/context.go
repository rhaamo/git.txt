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
)

// Context represents context of a request.
type Context struct {
	*macaron.Context
	Cache   cache.Cache
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	User	string // should be a models.User

	IsLogged    bool
	IsBasicAuth bool

}

// Title sets "Title" field in template data.
func (c *Context) Title(locale string) {
	c.Data["Title"] = c.Tr(locale)
}

// HTML responses template with given status.
func (ctx *Context) HTML(status int, name string) {
	log.Trace("Template: %s", name)
	ctx.Context.HTML(status, name)
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
		// ctx.User, ctx.IsBasicAuth = auth.SignedInUser(ctx.Context, ctx.Session)

		//if ctx.User != nil {
		//	ctx.IsLogged = true
		//	ctx.Data["IsLogged"] = ctx.IsLogged
		//	ctx.Data["LoggedUser"] = ctx.User
		//	//ctx.Data["LoggedUserID"] = ctx.User.ID
		//	//ctx.Data["LoggedUserName"] = ctx.User.Name
		//} else {
		//	ctx.Data["LoggedUserID"] = 0
		//	ctx.Data["LoggedUserName"] = ""
		//}

		ctx.Data["CSRFToken"] = x.GetToken()
		ctx.Data["CSRFTokenHTML"] = template.HTML(`<input type="hidden" name="_csrf" value="` + x.GetToken() + `">`)
		log.Trace("Session ID: %s", sess.ID())
		log.Trace("CSRF Token: %v", ctx.Data["CSRFToken"])

		c.Map(ctx)
	}
}