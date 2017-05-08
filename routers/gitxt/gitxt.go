package gitxt

import (
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/stuff/form"
	"dev.sigpipe.me/dashie/git.txt/stuff/sanitize"
	"fmt"
	"strings"
	"dev.sigpipe.me/dashie/git.txt/stuff/tool"
	"time"
	"dev.sigpipe.me/dashie/git.txt/stuff/repository"
	"os"
	log "gopkg.in/clog.v1"
	"gopkg.in/libgit2/git2go.v25"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"dev.sigpipe.me/dashie/git.txt/stuff/gite"
)

const (
	NEW = "gitxt/new"
	VIEW = "gitxt/view"
	LIST = "gitxt/list"
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
	u := &models.Gitxt{
		Hash: repositoryName,
		Description: f.Description,
		IsPrivate: !f.IsPublic,
	}

	if ctx.IsLogged {
		u.UserID = ctx.User.ID
		u.Anonymous = false
	} else {
		u.UserID = 0
		u.Anonymous = true
	}

	if err := models.CreateGitxt(u); err != nil {
		switch {
		case models.IsErrHashAlreadyExist(err):
			ctx.Data["Err_Hash"] = true
			ctx.RenderWithErr(ctx.Tr("gitxt_new.hash_been_taken"), NEW, &f)
		default:
			ctx.Handle(500, "NewPost", err)
		}
		return
	}

	// 5. Return render to gitxt view page

	log.Trace("Pushed repository %s to database as %i", repositoryName, u.ID)
	ctx.Redirect(setting.AppSubURL + "/" + repositoryUser + "/" + repositoryName)
}

// View gitxt
func View(ctx *context.Context) {
	ctx.Title("gitxt_view.title")
	ctx.PageIs("GitxtView")

	ctx.Data["repoDescription"] = ctx.Gitxt.Description
	ctx.Data["repoIsPrivate"] = ctx.Gitxt.IsPrivate
	ctx.Data["repoOwnerUsername"] = ctx.RepoOwnerUsername
	ctx.Data["repoHash"] = ctx.Gitxt.Hash

	// Get the files from git
	var repoSpec = "HEAD"

	repoPath := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Hash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Hash, err)
		ctx.Flash.Error("Error: could not open repository.")
		ctx.HTML(500, "GitxtView")
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty();
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Hash, err)
		ctx.Flash.Error("Error: repository is empty or corrupted")
		ctx.HTML(500, "GitxtView")
		return
	}

	// Get the repository tree
	repoTreeEntries, err := gite.GetWalkTreeWithContent(repo, "/")
	if err != nil {
		log.Warn("Cannot get repository tree entries %s: %s", ctx.Gitxt.Hash, err)
		ctx.Flash.Error("Error: cannot get repository tree entries")
		ctx.HTML(500, "GitxtView")
		return

	}

	ctx.Data["repoSpec"] = repoSpec
	ctx.Data["repoFiles"] = repoTreeEntries

	ctx.Success(VIEW)
}

func ListForUser(ctx *context.Context) {
	ctx.Title("gitxt_list.title")
	ctx.PageIs("GitxtList")

	ctx.Data["GitxtListIsUser"] = true

	ctx.Success(LIST)
}