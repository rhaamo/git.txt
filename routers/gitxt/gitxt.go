package gitxt

import (
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/stuff/form"
	"dev.sigpipe.me/dashie/git.txt/stuff/sanitize"
	// "dev.sigpipe.me/dashie/git.txt/models"
	"fmt"
	"strings"
	"dev.sigpipe.me/dashie/git.txt/stuff/tool"
	"time"
	"dev.sigpipe.me/dashie/git.txt/stuff/repository"
	"os"
	log "gopkg.in/clog.v1"
	"gopkg.in/libgit2/git2go.v25"
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
	// Craft a repository name from a SHA1 which should be pretty much unique
	repositoryName := tool.SHA1(time.Now().String() + ctx.Data["CSRFToken"].(string))

	var repositoryUser string
	if ctx.IsLogged {
		repositoryUser = ctx.Data["LoggedUserName"].(string)
	} else {
		repositoryUser = "anonymous"
	}

	// git init --blahblah bare
	repo, err := repository.InitRepository(repositoryUser, repositoryName)
	if err != nil {
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_repository_error")

		// Get repository path and remove it from filesystem
		RepoPath := repository.RepoPath(repositoryUser, repositoryName)
		if err := os.RemoveAll(RepoPath); err != nil {
			log.Warn("Cannot remove repository '%s': %s", RepoPath, err)
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_and_remove_repository_error")
		}

		log.Trace("Repository deleted: %s for %s", repositoryName, repositoryUser)

		ctx.Success(NEW)
		return
	}

	log.Trace("Repository created: %s for %s", repositoryName, repositoryUser)

	// 2. 3. Create the files and commit

	// Create the blobs objects
	// git whatever create blob
	var blobs []*git.Oid
	for i := range f.FilesFilename {
		blob, err := repo.CreateBlobFromBuffer([]byte(f.FilesContent[i]))
		if err != nil {
			log.Warn("init_error_create_blob: %s", err)
			ctx.Data["HasError"] = true
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_error_create_blob")
			ctx.Success(NEW)
			return
		}
		blobs = append(blobs, blob)
	}

	//
	repoIndex, err := repo.Index()
	if err != nil {
		log.Warn("init_error_get_index: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_error_get_index")
		ctx.Success(NEW)
		return
	}

	// Add the blobs to the index
	// git update-index --add --cacheinfo 100644 "$BLOB_ID" "myfile.txt"
	for i := range f.FilesFilename {
		indexEntry := &git.IndexEntry{
			Path: f.FilesFilename[i],
			Mode: git.FilemodeBlob,
			Id: blobs[i],
		}
		if repoIndex.Add(indexEntry) != nil {
			log.Warn("init_error_add_entry: %s", err)
			ctx.Data["HasError"] = true
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_error_add_entry")
			ctx.Success(NEW)
			return
		}
	}

	// Write the new tree
	// TREE_ID=$(git write-tree)
	repoTreeOid, err := repoIndex.WriteTree()
	if err != nil {
		log.Warn("init_error_index_write_tree: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_error_index_write_tree")
		ctx.Success(NEW)
		return

	}

	// Get latest tree
	repoTree, err := repo.LookupTree(repoTreeOid)
	if err != nil {
		log.Warn("init_error_lookup_tree: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_error_lookup_tree")
		ctx.Success(NEW)
		return

	}

	// NEW_COMMIT=$(echo "My commit message" | git commit-tree "$TREE_ID" -p "$PARENT_COMMIT")
	ruAuthor := &git.Signature{
		Name: repositoryUser,
		Email: "autocommit@git.txt",
	}
	_, err = repo.CreateCommit("HEAD", ruAuthor, ruAuthor, "First autocommit from git.txt", repoTree)
	if err != nil {
		log.Warn("init_error_commit: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.init_error_commit")
		ctx.Success(NEW)
		return

	}

	// git update-ref "refs/heads/$MY_BRANCH" "$NEW_COMMIT" "$PARENT_COMMIT"

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