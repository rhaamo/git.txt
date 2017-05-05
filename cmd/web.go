package cmd

import (
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/routers"
	"dev.sigpipe.me/dashie/git.txt/bindata"
	"path"
	"github.com/go-macaron/session"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/i18n"
	"strings"
	"fmt"
	log "gopkg.in/clog.v1"
	"net/http"
	"net"
	"net/http/fcgi"
	"os"
)

var Web = cli.Command{
	Name: "web",
	Usage: "Start web server",
	Description: "It starts a web server, great no ?",
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Server port"),
		stringFlag("config, c", "config/app.ini", "Custom config file path"),
	},
}

func newMacaron() *macaron.Macaron {
	m := macaron.New()
	if !setting.DisableRouterLog {
		m.Use(macaron.Logger())
	}

	m.Use(macaron.Recovery())

	if setting.Protocol == setting.SCHEME_FCGI {
		m.SetURLPrefix(setting.AppSubURL)
	}

	m.Use(macaron.Static(
		path.Join(setting.StaticRootPath, "public"),
		macaron.StaticOptions{
			SkipLogging: setting.DisableRouterLog,
		},
	))

	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:         path.Join(setting.StaticRootPath, "templates"),
		IndentJSON:        macaron.Env != macaron.PROD,
	}))

	localeNames, err := bindata.AssetDir("conf/locale")
	if err != nil {
		log.Fatal(4, "Fail to list locale files: %v", err)
	}
	localFiles := make(map[string][]byte)
	for _, name := range localeNames {
		localFiles[name] = bindata.MustAsset("conf/locale/" + name)
	}
	m.Use(i18n.I18n(i18n.Options{
		SubURL:          setting.AppSubURL,
		Files:           localFiles,
		Langs:           setting.Langs,
		Names:           setting.Names,
		DefaultLang:     "en-US",
		Redirect:        true,
	}))

	m.Use(cache.Cacher(cache.Options{
		Adapter:       setting.CacheAdapter,
		AdapterConfig: setting.CacheConn,
		Interval:      setting.CacheInterval,
	}))

	m.Use(session.Sessioner(setting.SessionConfig))

	m.Use(csrf.Csrfer(csrf.Options{
		Secret:     setting.SecretKey,
		Cookie:     setting.CSRFCookieName,
		SetCookie:  true,
		Header:     "X-Csrf-Token",
		CookiePath: setting.AppSubURL,
	}))

	// m.Use(toolbox.Toolboxer(m, toolbox.Options{
	// 	HealthCheckFuncs: []*toolbox.HealthCheckFuncDesc{
	// 		&toolbox.HealthCheckFuncDesc{
	// 			Desc: "Database connection",
	// 			Func: models.Ping,
	// 		},
	// 	},
	// }))

	m.Use(context.Contexter())

	return m

}

func runWeb(ctx *cli.Context) error {
	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
	}

	setting.InitConfig()

	m := newMacaron()

	m.Get("/", routers.Home)

	// robots.txt
	m.Get("/robots.txt", func(ctx *context.Context) {
		if setting.HasRobotsTxt {
			ctx.ServeFileContent(setting.RobotsTxtPath)
		} else {
			ctx.Error(404)
		}
	})

	// TODO not found handler
	// m.NotFound()

	if ctx.IsSet("port") {
		setting.AppURL = strings.Replace(setting.AppURL, setting.HTTPPort, ctx.String("port"), 1)
		setting.HTTPPort = ctx.String("port")
	}

	var listenAddr string
	if setting.Protocol == setting.SCHEME_UNIX_SOCKET {
		listenAddr = fmt.Sprintf("%s", setting.HTTPAddr)
	} else {
		listenAddr = fmt.Sprintf("%s:%s", setting.HTTPAddr, setting.HTTPPort)
	}
	log.Info("Listen: %v://%s%s", setting.Protocol, listenAddr, setting.AppSubURL)

	var err error
	switch setting.Protocol {
	case setting.SCHEME_HTTP:
		err = http.ListenAndServe(listenAddr, m)
	case setting.SCHEME_HTTPS:
		log.Fatal(2, "https not supported")
	case setting.SCHEME_FCGI:
		err = fcgi.Serve(nil, m)
	case setting.SCHEME_UNIX_SOCKET:
		os.Remove(listenAddr)

		var listener *net.UnixListener
		listener, err = net.ListenUnix("unix", &net.UnixAddr{listenAddr, "unix"})
		if err != nil {
			break // Handle error after switch
		}

		// FIXME: add proper implementation of signal capture on all protocols
		// execute this on SIGTERM or SIGINT: listener.Close()
		if err = os.Chmod(listenAddr, os.FileMode(setting.UnixSocketPermission)); err != nil {
			log.Fatal(4, "Failed to set permission of unix socket: %v", err)
		}
		err = http.Serve(listener, m)
	default:
		log.Fatal(4, "Invalid protocol: %s", setting.Protocol)
	}

	if err != nil {
		log.Fatal(4, "Fail to start server: %v", err)
	}

	return nil
}