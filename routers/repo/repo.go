// Copyright 2017 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style license

package repo

import (
	"github.com/go-macaron/macaron"
	"dev.sigpipe.me/dashie/git.txt/context"
	"strings"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/models/errors"
	"net/http"
	"dev.sigpipe.me/dashie/git.txt/stuff/tool"
	"dev.sigpipe.me/dashie/git.txt/setting"
	log "gopkg.in/clog.v1"
	"regexp"
	"time"
	"fmt"
	"path"
	"os"
	"github.com/Unknwon/com"
	"compress/gzip"
	"bytes"
	"os/exec"
	"strconv"
)

const (
	ENV_AUTH_USER_ID           = "GITXT_AUTH_USER_ID"
	ENV_AUTH_USER_NAME         = "GITXT_AUTH_USER_NAME"
	ENV_AUTH_USER_EMAIL        = "GITXT_AUTH_USER_EMAIL"
	ENV_REPO_OWNER_NAME        = "GITXT_REPO_OWNER_NAME"
	ENV_REPO_OWNER_SALT_MD5    = "GITXT_REPO_OWNER_SALT_MD5"
	ENV_REPO_ID                = "GITXT_REPO_ID"
	ENV_REPO_NAME              = "GITXT_REPO_NAME"
	ENV_REPO_CUSTOM_HOOKS_PATH = "GITXT_REPO_CUSTOM_HOOKS_PATH"
)

type HTTPContext struct {
	*context.Context
	OwnerName	string
	OwnerSalt	string
	RepoID		int64
	RepoName	string
	AuthUser	*models.User
}

// askCredentials responses HTTP header and status which informs client to provide credentials
func askCredentials(ctx *context.Context, status int, text string) {
	ctx.Resp.Header().Set("WWW-Authenticate", "Basic realm=\".\"")
	ctx.HandleText(status, text)
}

func HTTPContexter() macaron.Handler {
	return func(ctx *context.Context) {
		ownerName := ctx.Params(":user")
		repoName := strings.TrimSuffix(ctx.Params(":hash"), ".git")

		isPull := ctx.Query("service") == "git-upload-pack" || strings.HasSuffix(ctx.Req.URL.Path, "git-upload-pack") || ctx.Req.Method == "GET"

		owner, err := models.GetUserByName(ownerName)
		if err != nil && ownerName != "anonymous" {
			log.Trace("Could not found user: %s", ownerName)
			ctx.NotFoundOrServerError("GetUserByName", errors.IsUserNotExist, err)
			return
		}

		repo, err := models.GetRepositoryByName(ownerName, repoName)
		if err != nil {
			log.Trace("Could not found repository: %s", repoName)
			ctx.NotFoundOrServerError("GetRepositoryByName", errors.IsRepoNotExist, err)
			return
		}

		if isPull {
			ctx.Map(&HTTPContext{
				Context: ctx,
			})
			log.Trace("It's a pull request")
			return
		}

		// In case user requested a wrong URL and not intended to access Git objects.
		action := ctx.Params("*")
		if !strings.Contains(action, "git-") &&
			!strings.Contains(action, "info/") &&
			!strings.Contains(action, "HEAD") &&
			!strings.Contains(action, " objects/") {
			log.Trace("Whoops could not match anything from this action: %s", action)
			ctx.NotFound()
			return
		}

		// Handle HTTP Basic Auth
		authHead := ctx.Req.Header.Get("Authorization")
		if len(authHead) == 0 {
			askCredentials(ctx, http.StatusUnauthorized, "")
		}

		auths := strings.Fields(authHead)
		if len(auths) != 2 || auths[0] != "Basic" {
			askCredentials(ctx, http.StatusUnauthorized, "")
			return
		}

		authUsername, authPassword, err := tool.BasicAuthDecode(auths[1])
		if err != nil {
			askCredentials(ctx, http.StatusUnauthorized, "")
			return
		}

		authUser, err := models.UserLogin(authUsername, authPassword)
		log.Trace("%s", err)
		if err != nil {
			if errors.IsUserNotExist(err) {
				askCredentials(ctx, http.StatusUnauthorized, "")
			} else {
				ctx.Handle(http.StatusInternalServerError, "UserLogin", err)
			}
			return
		}

		log.Trace("HTTPGit - Authenticated user: %s", authUser.UserName)

		// Reject if not pulling and user doesn't match repo ID
		if authUser.ID != repo.UserID  && !isPull {
			askCredentials(ctx, http.StatusForbidden, "User permission denied")
			return
		}

		ctx.Map(&HTTPContext{
			Context: ctx,
			OwnerName: ownerName,
			OwnerSalt: owner.Salt,
			RepoID: repo.ID,
			RepoName: repoName,
			AuthUser: authUser,
		})

	}
}

type serviceHandler struct {
	w    http.ResponseWriter
	r    *http.Request
	dir  string
	file string

	authUser  *models.User
	ownerName string
	ownerSalt string
	repoID    int64
	repoName  string
}

func (h *serviceHandler) setHeaderNoCache() {
	h.w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	h.w.Header().Set("Pragma", "no-cache")
	h.w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func (h *serviceHandler) setHeaderCacheForever() {
	now := time.Now().Unix()
	expires := now + 31536000
	h.w.Header().Set("Date", fmt.Sprintf("%d", now))
	h.w.Header().Set("Expires", fmt.Sprintf("%d", expires))
	h.w.Header().Set("Cache-Control", "public, max-age=31536000")
}

func (h *serviceHandler) sendFile(contentType string) {
	reqFile := path.Join(h.dir, h.file)
	fi, err := os.Stat(reqFile)
	if os.IsNotExist(err) {
		h.w.WriteHeader(http.StatusNotFound)
		return
	}

	h.w.Header().Set("Content-Type", contentType)
	h.w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))
	h.w.Header().Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))
	http.ServeFile(h.w, h.r, reqFile)
}

type ComposeHookEnvsOptions struct {
	AuthUser  *models.User
	OwnerName string
	OwnerSalt string
	RepoID    int64
	RepoName  string
	RepoPath  string
}

func ComposeHookEnvs(opts ComposeHookEnvsOptions) []string {
	envs := []string{
		"SSH_ORIGINAL_COMMAND=1",
		ENV_AUTH_USER_ID + "=" + com.ToStr(opts.AuthUser.ID),
		ENV_AUTH_USER_NAME + "=" + opts.AuthUser.UserName,
		ENV_AUTH_USER_EMAIL + "=" + opts.AuthUser.Email,
		ENV_REPO_OWNER_NAME + "=" + opts.OwnerName,
		ENV_REPO_OWNER_SALT_MD5 + "=" + tool.MD5(opts.OwnerSalt),
		ENV_REPO_ID + "=" + com.ToStr(opts.RepoID),
		ENV_REPO_NAME + "=" + opts.RepoName,
		ENV_REPO_CUSTOM_HOOKS_PATH + "=" + path.Join(opts.RepoPath, "custom_hooks"),
	}
	return envs
}

func serviceRPC(h serviceHandler, service string) {
	defer h.r.Body.Close()

	if h.r.Header.Get("Content-Type") != fmt.Sprintf("application/x-git-%s-request", service) {
		h.w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", service))

	var (
		reqBody = h.r.Body
		err     error
	)

	// Handle GZIP
	if h.r.Header.Get("Content-Encoding") == "gzip" {
		reqBody, err = gzip.NewReader(reqBody)
		if err != nil {
			log.Error(2, "HTTP.Get: fail to create gzip reader: %v", err)
			h.w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var stderr bytes.Buffer
	cmd := exec.Command(setting.GitBinary, service, "--stateless-rpc", h.dir)
	if service == "receive-pack" {
		cmd.Env = append(os.Environ(), ComposeHookEnvs(ComposeHookEnvsOptions{
			AuthUser:  h.authUser,
			OwnerName: h.ownerName,
			OwnerSalt: h.ownerSalt,
			RepoID:    h.repoID,
			RepoName:  h.repoName,
			RepoPath:  h.dir,
		})...)
	}
	cmd.Dir = h.dir
	cmd.Stdout = h.w
	cmd.Stderr = &stderr
	cmd.Stdin = reqBody
	if err = cmd.Run(); err != nil {
		log.Error(2, "HTTP.serviceRPC: fail to serve RPC '%s': %v - %s", service, err, stderr)
		h.w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func serviceUploadPack(h serviceHandler) {
	serviceRPC(h, "upload-pack")
}

func serviceReceivePack(h serviceHandler) {
	serviceRPC(h, "receive-pack")
}

func getServiceType(r *http.Request) string {
	serviceType := r.FormValue("service")
	if !strings.HasPrefix(serviceType, "git-") {
		return ""
	}
	return strings.TrimPrefix(serviceType, "git-")
}

// FIXME: use process module
func gitCommand(dir string, args ...string) []byte {
	cmd := exec.Command(setting.GitBinary, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		log.Error(2, fmt.Sprintf("Git: %v - %s", err, out))
	}
	return out
}

func updateServerInfo(dir string) []byte {
	return gitCommand(dir, "update-server-info")
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)
	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}

func getInfoRefs(h serviceHandler) {
	h.setHeaderNoCache()
	service := getServiceType(h.r)
	if service != "upload-pack" && service != "receive-pack" {
		updateServerInfo(h.dir)
		h.sendFile("text/plain; charset=utf-8")
		return
	}

	refs := gitCommand(h.dir, service, "--stateless-rpc", "--advertise-refs", ".")
	h.w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service))
	h.w.WriteHeader(http.StatusOK)
	h.w.Write(packetWrite("# service=git-" + service + "\n"))
	h.w.Write([]byte("0000"))
	h.w.Write(refs)
}

func getTextFile(h serviceHandler) {
	h.setHeaderNoCache()
	h.sendFile("text/plain")
}

func getInfoPacks(h serviceHandler) {
	h.setHeaderCacheForever()
	h.sendFile("text/plain; charset=utf-8")
}

func getLooseObject(h serviceHandler) {
	h.setHeaderCacheForever()
	h.sendFile("application/x-git-loose-object")
}

func getPackFile(h serviceHandler) {
	h.setHeaderCacheForever()
	h.sendFile("application/x-git-packed-objects")
}

func getIdxFile(h serviceHandler) {
	h.setHeaderCacheForever()
	h.sendFile("application/x-git-packed-objects-toc")
}

var routes = []struct {
	reg     *regexp.Regexp
	method  string
	handler func(serviceHandler)
}{
	{regexp.MustCompile("(.*?)/git-upload-pack$"), "POST", serviceUploadPack},
	{regexp.MustCompile("(.*?)/git-receive-pack$"), "POST", serviceReceivePack},
	{regexp.MustCompile("(.*?)/info/refs$"), "GET", getInfoRefs},
	{regexp.MustCompile("(.*?)/HEAD$"), "GET", getTextFile},
	{regexp.MustCompile("(.*?)/objects/info/alternates$"), "GET", getTextFile},
	{regexp.MustCompile("(.*?)/objects/info/http-alternates$"), "GET", getTextFile},
	{regexp.MustCompile("(.*?)/objects/info/packs$"), "GET", getInfoPacks},
	{regexp.MustCompile("(.*?)/objects/info/[^/]*$"), "GET", getTextFile},
	{regexp.MustCompile("(.*?)/objects/[0-9a-f]{2}/[0-9a-f]{38}$"), "GET", getLooseObject},
	{regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.pack$"), "GET", getPackFile},
	{regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.idx$"), "GET", getIdxFile},
}

func getGitRepoPath(repoDir string) (string, error) {
	if !strings.HasSuffix(repoDir, ".git") {
		repoDir += ".git"
	}

	filename := path.Join(setting.RepositoryRoot, repoDir)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", err
	}

	return filename, nil
}

func HTTP(ctx *HTTPContext) {
	for _, route := range routes {
		reqPath := strings.ToLower(ctx.Req.URL.Path)
		m := route.reg.FindStringSubmatch(reqPath)
		if m == nil {
			continue
		}

		// We perform check here because routes matched in cmd/web.go is wider than needed,
		// but we only want to output this message only if user is really trying to access
		// Git HTTP endpoints.
		if setting.DisableHttpGit {
			ctx.HandleText(http.StatusForbidden, "Interacting with repositories by HTTP protocol is not disabled")
			return
		}

		if route.method != ctx.Req.Method {
			ctx.NotFound()
			return
		}

		file := strings.TrimPrefix(reqPath, m[1]+"/")
		dir, err := getGitRepoPath(m[1])
		if err != nil {
			log.Warn("HTTP.getGitRepoPath: %v", err)
			ctx.NotFound()
			return
		}

		route.handler(serviceHandler{
			w:    ctx.Resp,
			r:    ctx.Req.Request,
			dir:  dir,
			file: file,

			authUser:  ctx.AuthUser,
			ownerName: ctx.OwnerName,
			ownerSalt: ctx.OwnerSalt,
			repoID:    ctx.RepoID,
			repoName:  ctx.RepoName,
		})
		return
	}

	log.Trace("Not found %s", ctx.Req.Request)
	ctx.NotFound()
}