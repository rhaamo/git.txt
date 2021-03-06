package template

import (
	"container/list"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"dev.sigpipe.me/dashie/git.txt/stuff/markup"
	"dev.sigpipe.me/dashie/git.txt/stuff/template/highlight"
	"dev.sigpipe.me/dashie/git.txt/stuff/tool"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"html/template"
	"mime"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// NewFuncMap initialize the Functions Map
func NewFuncMap() []template.FuncMap {
	return []template.FuncMap{map[string]interface{}{
		"GoVer": func() string {
			return strings.Title(runtime.Version())
		},
		"AppName": func() string {
			return setting.AppName
		},
		"AppSubURL": func() string {
			return setting.AppSubURL
		},
		"AppURL": func() string {
			return setting.AppURL
		},
		"AppVer": func() string {
			if len(setting.BuildGitHash) > 0 {
				return fmt.Sprintf("%s - %s", setting.AppVer, setting.BuildGitHash)
			}
			return setting.AppVer
		},
		"AppDomain": func() string {
			return setting.Domain
		},
		"CanRegister": func() bool {
			return setting.CanRegister
		},
		"AnonymousCreate": func() bool {
			return setting.AnonymousCreate
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		"FileSize": tool.FileSize,
		"Safe":     safe,
		"Sanitize": bluemonday.UGCPolicy().Sanitize,
		"Str2html": str2html,
		"Add": func(a, b int) int {
			return a + b
		},
		"DateFmtLong": func(t time.Time) string {
			return t.Format(time.RFC1123Z)
		},
		"DateFmtShort": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
		"List": listWhatever,
		"SubStr": func(str string, start, length int) string {
			if len(str) == 0 {
				return ""
			}
			end := start + length
			if length == -1 {
				end = len(str)
			}
			if len(str) < end {
				return str
			}
			return str[start:end]
		},
		"Join":        strings.Join,
		"Sha1":        sha1,
		"ShortSHA1":   tool.ShortSHA1,
		"MD5":         tool.MD5,
		"EscapePound": escapePound,
		"FilenameIsImage": func(filename string) bool {
			mimeType := mime.TypeByExtension(filepath.Ext(filename))
			return strings.HasPrefix(mimeType, "image/")
		},
		"FilenameToHighlightClass": func(filename string) string {
			return highlight.FileNameToHighlightClass(filename)
		},
		"IsMarkdown": func(filename string) bool {
			return markup.IsMarkdownFile(filename)
		},
		"ToMarkdown": func(content string) string {
			return string(markup.Markdown(content, setting.AppSubURL)[:])
		},
		"IsPdf": func(mime string) bool {
			return strings.EqualFold(mime, "application/pdf")
		},
		"CounterGitxt": func() int64 {
			c, _ := models.GetCounterGitxts()
			return c.Count
		},
		"CounterGitxtManaged": func() int64 {
			c, _ := models.GetCounterGitxtsManaged()
			return c.Count
		},
	}}
}

func safe(raw string) template.HTML {
	return template.HTML(raw)
}

func listWhatever(l *list.List) chan interface{} {
	e := l.Front()
	c := make(chan interface{})
	go func() {
		for e != nil {
			c <- e.Value
			e = e.Next()
		}
		close(c)
	}()
	return c
}

func sha1(str string) string {
	return tool.SHA1(str)
}

func escapePound(str string) string {
	return strings.NewReplacer("%", "%25", "#", "%23", " ", "%20", "?", "%3F").Replace(str)
}

func str2html(raw string) template.HTML {
	return template.HTML(markup.Sanitize(raw))
}
