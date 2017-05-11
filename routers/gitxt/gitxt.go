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
	"path/filepath"
)

const (
	NEW = "gitxt/new"
	VIEW = "gitxt/view"
	LIST = "gitxt/list"
	EDIT = "gitxt/edit"
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
	// Reject-redirect if not logged-in and anonymous create is deactivated
	if !setting.AnonymousCreate && !ctx.IsLogged {
		ctx.Redirect(setting.AppSubURL + "/")
		return
	}

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
		if len(filepath.Ext(f.FilesFilename[i])) == 0 {
			// No extension, forces .txt
			f.FilesFilename[i] = fmt.Sprintf("%s.txt", f.FilesFilename[i])
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
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_init_repository")

		// Get repository path and remove it from filesystem
		RepoPath := repository.RepoPath(repositoryUser, repositoryName)
		if err := os.RemoveAll(RepoPath); err != nil {
			log.Warn("Cannot remove repository '%s': %s", RepoPath, err)
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_remove_repository")
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
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_create_blob")
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
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_index")
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
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_add_entry")
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
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_index_write_tree")
		ctx.Success(NEW)
		return

	}

	// Get latest tree
	repoTree, err := repo.LookupTree(repoTreeOid)
	if err != nil {
		log.Warn("init_error_lookup_tree: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_lookup_tree")
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
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_commit")
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

	ctx.Data["repoDescription"] = ctx.Gitxt.Gitxt.Description
	ctx.Data["repoIsPrivate"] = ctx.Gitxt.Gitxt.IsPrivate
	ctx.Data["repoOwnerUsername"] = ctx.RepoOwnerUsername
	ctx.Data["repoHash"] = ctx.Gitxt.Gitxt.Hash

	// Get the files from git
	var repoSpec = "HEAD"

	repoPath := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Gitxt.Hash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: could not open repository.")
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty();
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: repository is empty or corrupted")
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Get the repository tree
	repoTreeEntries, err := gite.GetWalkTreeWithContent(repo, "/")
	if err != nil {
		log.Warn("Cannot get repository tree entries %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: cannot get repository tree entries")
		ctx.Handle(500, "GitxtView", err)
		return

	}

	ctx.Data["repoSpec"] = repoSpec
	ctx.Data["repoFiles"] = repoTreeEntries
	ctx.Data["IsOwner"] = ctx.Gitxt.Owner

	ctx.Success(VIEW)
}

// List uploads, manage auth'ed user or not and from /:user too
func ListUploads(ctx *context.Context) {
	ctx.Title("gitxt_list.title")
	ctx.PageIs("GitxtList")

	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}
	ctx.Data["PageNumber"] = page

	opts := &models.GitxtOptions{
		PageSize: 10,	// TODO: put this in config
		Page: page,
	}

	if ctx.Data["GetAll"] != nil {
		opts.GetAll = true
	}

	if ctx.RepoOwnerUsername != "" {
		ctx.Data["GitxtListIsUser"] = true
		opts.UserID = ctx.Gitxt.User.ID
	} else {
		ctx.Data["GitxtListIsUser"] = false
	}
	ctx.Data["RepoOwnerUsername"] = ctx.RepoOwnerUsername

	if ctx.IsLogged {
		opts.WithPrivate = true
	} else {
		opts.WithPrivate = false
	}

	listOfGitxts, gitxtsCount, err := models.GetGitxts(opts)
	if err != nil {
		log.Warn("Cannot get Gitxts with opts %v, %s", opts, err)
		ctx.Flash.Error("Error while getting list of Gitxts")
		ctx.Handle(500, "ListUploads", err)
		return
	}

	ctx.Data["gitxts"] = listOfGitxts
	ctx.Data["gitxts_count"] = gitxtsCount

	ctx.Success(LIST)
}

func DeletePost(ctx *context.Context, f form.GitxtDelete) {
	if ctx.HasError() {
		ctx.JSONSuccess(map[string]interface{}{
			"error": ctx.Data["ErrorMsg"],
			"redirect": false,
		})
		return
	}

	if ctx.Data["LoggedUserID"] != ctx.Gitxt.Gitxt.UserID {
		ctx.JSONSuccess(map[string]interface{}{
			"error": "Unauthorized",
			"redirect": false,
		})
	}
}

func Edit(ctx *context.Context) {
	ctx.Title("gitxt_edit.title")
	ctx.PageIs("GitxtEdit")

	if ctx.Data["LoggedUserID"] != ctx.Gitxt.Gitxt.UserID {
		ctx.Flash.Error("Unauthorized")
		ctx.Redirect(setting.AppSubURL + ctx.RepoOwnerUsername + "/" + ctx.Gitxt.Gitxt.Hash)
		return
	}

	ctx.Data["description"] = ctx.Gitxt.Gitxt.Description
	ctx.Data["repoIsPrivate"] = ctx.Gitxt.Gitxt.IsPrivate
	ctx.Data["repoOwnerUsername"] = ctx.RepoOwnerUsername
	ctx.Data["repoHash"] = ctx.Gitxt.Gitxt.Hash

	// Get the files from git
	var repoSpec = "HEAD"

	repoPath := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Gitxt.Hash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: could not open repository.")
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty();
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: repository is empty or corrupted")
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Get the repository tree
	repoTreeEntries, err := gite.GetWalkTreeWithContent(repo, "/")
	if err != nil {
		log.Warn("Cannot get repository tree entries %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error("Error: cannot get repository tree entries")
		ctx.Handle(500, "GitxtView", err)
		return

	}

	ctx.Data["repoSpec"] = repoSpec

	var FilesContent []string
	var FilesFilename []string

	for i := range repoTreeEntries {
		FilesContent = append(FilesContent, repoTreeEntries[i].Content)
		FilesFilename = append(FilesFilename, repoTreeEntries[i].Path)
	}

	ctx.Data["FilesContent"] = FilesContent
	ctx.Data["FilesFilename"] = FilesFilename

	ctx.Success(EDIT)
}

func EditPost(ctx *context.Context, f form.GitxtEdit) {
	if !ctx.IsLogged {
		ctx.Redirect(setting.AppSubURL + "/")
		return
	}

	if ctx.Data["LoggedUserID"] != ctx.Gitxt.Gitxt.UserID {
		ctx.Flash.Error("Unauthorized")
		ctx.Redirect(setting.AppSubURL + ctx.RepoOwnerUsername + "/" + ctx.Gitxt.Gitxt.Hash)
		return
	}

	ctx.Title("gitxt_edit.title")
	ctx.PageIs("GitxtEditPost")

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
		if len(filepath.Ext(f.FilesFilename[i])) == 0 {
			// No extension, forces .txt
			f.FilesFilename[i] = fmt.Sprintf("%s.txt", f.FilesFilename[i])
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

	repositoryUser := ctx.Gitxt.User.UserName
	repositoryHash := ctx.Gitxt.Gitxt.Hash

	repoPath := repository.RepoPath(repositoryUser, repositoryHash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.could_not_open"))
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty();
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.repo_corrupt_or_empty"))
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	// git whatever create blob
	var blobs []*git.Oid
	for i := range f.FilesFilename {
		blob, err := repo.CreateBlobFromBuffer([]byte(f.FilesContent[i]))
		if err != nil {
			log.Warn("init_error_create_blob: %s", err)
			ctx.Data["HasError"] = true
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_create_blob")
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
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_index")
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
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_add_entry")
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
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_index_write_tree")
		ctx.Success(NEW)
		return

	}

	// Get latest tree
	repoTree, err := repo.LookupTree(repoTreeOid)
	if err != nil {
		log.Warn("init_error_lookup_tree: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_lookup_tree")
		ctx.Success(NEW)
		return

	}

	// Get repo head
	repoHead, err := repo.Head()
	if err != nil {
		log.Warn("git_error_get_head: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_head")
		ctx.Success(NEW)
		return
	}

	// Get latest commit
	headCommit, err := repo.LookupCommit(repoHead.Target())
	if err != nil {
		log.Warn("git_error_get_head_commit: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_head_commit")
		ctx.Success(NEW)
		return
	}

	// NEW_COMMIT=$(echo "My commit message" | git commit-tree "$TREE_ID" -p "$PARENT_COMMIT")
	ruAuthor := &git.Signature{
		Name: repositoryUser,
		Email: "autocommit@git.txt",
	}
	_, err = repo.CreateCommit("HEAD", ruAuthor, ruAuthor, "Autocommit from git.txt", repoTree, headCommit)
	if err != nil {
		log.Warn("init_error_commit: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_commit")
		ctx.Success(NEW)
		return

	}

	// git update-ref "refs/heads/$MY_BRANCH" "$NEW_COMMIT" "$PARENT_COMMIT"

	// 4. Insert info in database
	ctx.Gitxt.Gitxt.Description = f.Description
	if err := models.UpdateGitxt(ctx.Gitxt.Gitxt); err != nil {
		switch {
		default:
			ctx.Handle(500, "EditPost", err)
		}
		return
	}

	// 5. Return render to gitxt view page

	log.Trace("Edit Pushed repository %s - %i", ctx.Gitxt.Gitxt.Hash, ctx.Gitxt.Gitxt.ID)
	ctx.Redirect(setting.AppSubURL + "/" + repositoryUser + "/" + ctx.Gitxt.Gitxt.Hash)

}
