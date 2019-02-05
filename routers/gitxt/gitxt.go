package gitxt

import (
	"bytes"
	"dev.sigpipe.me/dashie/git.txt/context"
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"dev.sigpipe.me/dashie/git.txt/stuff/form"
	"dev.sigpipe.me/dashie/git.txt/stuff/gite"
	"dev.sigpipe.me/dashie/git.txt/stuff/markup"
	"dev.sigpipe.me/dashie/git.txt/stuff/repository"
	"dev.sigpipe.me/dashie/git.txt/stuff/sanitize"
	"dev.sigpipe.me/dashie/git.txt/stuff/tool"
	"fmt"
	"github.com/Unknwon/paginater"
	log "gopkg.in/clog.v1"
	"gopkg.in/libgit2/git2go.v27"
	gotemplate "html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	tmplNew  = "gitxt/new"
	tmplView = "gitxt/view"
	tmplList = "gitxt/list"
	tmplEdit = "gitxt/edit"
)

// New GET
func New(ctx *context.Context) {
	ctx.Title("gitxt_new.title")
	ctx.PageIs("GitxtNew")

	// One initial file
	ctx.Data["FilesContent"] = []string{""}
	ctx.Data["FilesFilename"] = []string{""}

	// Initial expiry
	ctx.Data["ExpiryHours"] = 0

	ctx.Success(tmplNew)
}

// NewPost POST
func NewPost(ctx *context.Context, f form.Gitxt) {
	// Reject-redirect if not logged-in and anonymous create is deactivated
	if !setting.AnonymousCreate && !ctx.IsLogged {
		ctx.Redirect(setting.AppSubURL + "/")
		return
	}

	ctx.Title("gitxt_new.title")
	ctx.PageIs("GitxtNewPost")

	for i := range f.FilesFilename {
		// For each filename sanitize it
		f.FilesFilename[i] = sanitize.Filename(f.FilesFilename[i])
		f.FilesFilename[i] = sanitize.RemoveFirstSlash(f.FilesFilename[i])
		if len(f.FilesFilename[i]) == 0 || f.FilesFilename[i] == "." {
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
	ctx.Data["ExpiryHours"] = f.ExpiryHours

	// We got an error in the manual validation step, render with error
	if ctx.HasError() {
		ctx.Success(tmplNew)
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

		ctx.Success(tmplNew)
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
			ctx.Success(tmplNew)
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
		ctx.Success(tmplNew)
		return
	}

	// Add the blobs to the index
	// git update-index --add --cacheinfo 100644 "$BLOB_ID" "myfile.txt"
	for i := range f.FilesFilename {
		indexEntry := &git.IndexEntry{
			Path: f.FilesFilename[i],
			Mode: git.FilemodeBlob,
			Id:   blobs[i],
		}
		if repoIndex.Add(indexEntry) != nil {
			log.Warn("init_error_add_entry: %s", err)
			ctx.Data["HasError"] = true
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_add_entry")
			ctx.Flash.Error(ctx.Tr("gitxt_git.error_add_entry"), true)
			ctx.Success(tmplNew)
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
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_index_write_tree"), true)
		ctx.Success(tmplNew)
		return

	}

	// Get latest tree
	repoTree, err := repo.LookupTree(repoTreeOid)
	if err != nil {
		log.Warn("init_error_lookup_tree: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_lookup_tree")
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_lookup_tree"), true)
		ctx.Success(tmplNew)
		return

	}

	// NEW_COMMIT=$(echo "My commit message" | git commit-tree "$TREE_ID" -p "$PARENT_COMMIT")
	ruAuthor := &git.Signature{
		Name:  repositoryUser,
		Email: "autocommit@git.txt",
	}
	_, err = repo.CreateCommit("HEAD", ruAuthor, ruAuthor, "First autocommit from git.txt", repoTree)
	if err != nil {
		log.Warn("init_error_commit: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_commit")
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_commit"), true)
		ctx.Success(tmplNew)
		return

	}

	// git update-ref "refs/heads/$MY_BRANCH" "$NEW_COMMIT" "$PARENT_COMMIT"

	// 4. Insert info in database
	u := &models.Gitxt{
		Hash:        repositoryName,
		Description: f.Description,
		IsPrivate:   !f.IsPublic,
		ExpiryHours: f.ExpiryHours,
		ExpiryUnix:  time.Now().Add(time.Hour * time.Duration(f.ExpiryHours)).Unix(),
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
			ctx.RenderWithErr(ctx.Tr("gitxt_new.hash_been_taken"), tmplNew, &f)
		default:
			ctx.Handle(500, "NewPost", err)
		}
		return
	}

	counter, _ := models.GetCounterGitxts()
	counterManaged, _ := models.GetCounterGitxtsManaged()
	models.UpdateCounterGitxts(counter.Count + 1)
	models.UpdateCounterGitxtsManaged(counterManaged.Count + 1)

	// 5. Return render to gitxt view page
	log.Trace("Pushed repository %s to database as %d", repositoryName, u.ID)
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
	ctx.Data["expiry"] = ctx.Gitxt.Gitxt.ExpiryHours
	ctx.Data["expiryOn"] = ctx.Gitxt.Gitxt.Expiry
	ctx.Data["createdOn"] = ctx.Gitxt.Gitxt.Created

	// Get the files from git
	var repoSpec = "HEAD"

	repoPath := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Gitxt.Hash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.could_not_open"), true)
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty()
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.repo_corrupt_or_empty"), true)
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Get the repository tree
	repoTreeEntries, err := gite.GetWalkTreeWithContent(repo, "/")
	if err != nil {
		log.Warn("Cannot get repository tree entries %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.cannot_get_tree_entries"), true)
		ctx.Handle(500, "GitxtView", err)
		return

	}

	for idx := range repoTreeEntries {
		if repoTreeEntries[idx].IsBinary || markup.IsMarkdownFile(repoTreeEntries[idx].Path) {
			continue
		}

		var output bytes.Buffer
		lines := strings.Split(repoTreeEntries[idx].Content, "\n")
		for index, line := range lines {
			output.WriteString(fmt.Sprintf(`<li class="L%d" rel="L%d">%s</li>`, index+1, index+1, gotemplate.HTMLEscapeString(line)) + "\n")
		}
		// Reset the Content and set the ContentH(tml)
		repoTreeEntries[idx].Content = ""
		repoTreeEntries[idx].ContentH = gotemplate.HTML(output.String())

		output.Reset()
		for i := 0; i < len(lines); i++ {
			output.WriteString(fmt.Sprintf(`<span id="L%d">%d</span>`, i+1, i+1))
		}
		repoTreeEntries[idx].LineNos = gotemplate.HTML(output.String())
	}

	ctx.Data["repoSpec"] = repoSpec
	ctx.Data["repoFiles"] = repoTreeEntries
	if ctx.IsLogged == true && !ctx.Gitxt.Gitxt.Anonymous {
		ctx.Data["IsOwner"] = ctx.Gitxt.User.ID == ctx.User.ID
	} else {
		ctx.Data["IsOwner"] = false
	}

	ctx.Success(tmplView)
}

// RawFile GET
func RawFile(ctx *context.Context) {
	file := ctx.Params("*")

	repoPath := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Gitxt.Hash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		ctx.ServerError("Cannot open repository", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty()
	if err != nil || isEmpty {
		ctx.ServerError("Repository empty or corrupted", err)
		return
	}

	// Get blob fileEntry struct
	treeFile, err := gite.GetTreeFileNoLimit(repo, file)
	if err != nil {
		ctx.ServerError("Cannot get file", err)
		return
	}

	// File not found
	if treeFile.ID == "" {
		ctx.Error(404)
		return
	}

	if treeFile.Size > setting.Bloby.MaxRawSize {
		ctx.Handle(500, "MaxRawSize", nil)
		return
	}

	// download if isBinary
	if treeFile.IsBinary && !(treeFile.MimeType == "image/png") && !(treeFile.MimeType == "application/pdf") {
		ctx.ServeContent(file, bytes.NewReader(treeFile.ContentB))
	} else {
		if strings.HasPrefix(treeFile.MimeType, "text/") {
			ctx.ServeContentNoDownload(file, "text/plain", bytes.NewReader(treeFile.ContentB))
		} else {
			ctx.ServeContentNoDownload(file, treeFile.MimeType, bytes.NewReader(treeFile.ContentB))
		}
	}
}

// ListUploads manage auth'ed user or not and from /:user too
func ListUploads(ctx *context.Context) {
	ctx.Title("gitxt_list.title")
	ctx.PageIs("GitxtList")

	page := ctx.QueryInt("page")
	if page <= 0 {
		page = 1
	}
	ctx.Data["PageNumber"] = page

	opts := &models.GitxtOptions{
		PageSize: 10, // TODO: put this in config
		Page:     page,
	}

	if ctx.Data["GetAll"] != nil {
		opts.GetAll = true
	}

	if ctx.RepoOwnerUsername != "" {
		ctx.Data["GitxtListIsUser"] = true
		if ctx.RepoOwnerUsername == "anonymous" {
			opts.UserID = 0
		} else {
			opts.UserID = ctx.Gitxt.User.ID
		}
		ctx.PageIs("GitxtListUser")
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
		ctx.Flash.Error(ctx.Tr("gitxt_list.error_getting_list"), true)
		ctx.Handle(500, "ListUploads", err)
		return
	}

	ctx.Data["gitxts"] = listOfGitxts
	ctx.Data["gitxts_count"] = gitxtsCount

	ctx.Data["Total"] = gitxtsCount
	ctx.Data["Page"] = paginater.New(int(gitxtsCount), opts.PageSize, page, 5)

	ctx.Success(tmplList)
}

// DeletePost GET
func DeletePost(ctx *context.Context, f form.GitxtDelete) {
	if ctx.HasError() {
		ctx.JSONSuccess(map[string]interface{}{
			"error":    ctx.Data["ErrorMsg"],
			"redirect": false,
		})
		return
	}

	if ctx.Data["LoggedUserID"] != ctx.Gitxt.Gitxt.UserID {
		ctx.JSONSuccess(map[string]interface{}{
			"error":    ctx.Tr("user.unauthorized"),
			"redirect": false,
		})
	}

	err := models.DeleteRepository(ctx.Gitxt.User.ID, ctx.Gitxt.Gitxt.ID)
	if err != nil {
		ctx.Flash.Error(ctx.Tr("gitxt_delete.error_deleting"), true)
		log.Warn("DeletePost.DeleteRepository: %v", err)
		ctx.JSONSuccess(map[string]interface{}{
			"error":    ctx.Tr("gitxt_delete.error_deleting"),
			"redirect": false,
		})
		return
	}

	counter, _ := models.GetCounterGitxts()
	if counter.Count <= 0 {
		// This should not happens but well, anyway
		models.UpdateCounterGitxts(counter.Count)
	} else {
		models.UpdateCounterGitxts(counter.Count - 1)
	}

	ctx.JSONSuccess(map[string]interface{}{
		"error":    nil,
		"redirect": setting.AppSubURL + "/",
	})
	return
}

// Edit git.txt
func Edit(ctx *context.Context) {
	ctx.Title("gitxt_edit.title")
	ctx.PageIs("GitxtEdit")

	if ctx.Data["LoggedUserID"] != ctx.Gitxt.Gitxt.UserID {
		ctx.Flash.Error(ctx.Tr("user.unauthorized"))
		ctx.Redirect(setting.AppSubURL + "/" + ctx.RepoOwnerUsername + "/" + ctx.Gitxt.Gitxt.Hash)
		return
	}

	ctx.Data["description"] = ctx.Gitxt.Gitxt.Description
	ctx.Data["repoIsPrivate"] = ctx.Gitxt.Gitxt.IsPrivate
	ctx.Data["repoOwnerUsername"] = ctx.RepoOwnerUsername
	ctx.Data["repoHash"] = ctx.Gitxt.Gitxt.Hash
	ctx.Data["ExpiryHours"] = ctx.Gitxt.Gitxt.ExpiryHours

	// Get the files from git
	var repoSpec = "HEAD"

	repoPath := repository.RepoPath(ctx.RepoOwnerUsername, ctx.Gitxt.Gitxt.Hash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.could_not_open"), true)
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty()
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.repo_corrupt_or_empty"), true)
		ctx.Handle(500, "GitxtView", err)
		return
	}

	// Get the repository tree
	repoTreeEntries, err := gite.GetWalkTreeWithContent(repo, "/")
	if err != nil {
		log.Warn("Cannot get repository tree entries %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.cannot_get_tree_entries"), true)
		ctx.Handle(500, "GitxtView", err)
		return

	}

	ctx.Data["repoSpec"] = repoSpec

	var FilesContent []string
	var FilesFilename []string
	var FilesNotHandled []bool

	for i := range repoTreeEntries {
		FilesContent = append(FilesContent, repoTreeEntries[i].Content)
		FilesFilename = append(FilesFilename, repoTreeEntries[i].Path)
		if repoTreeEntries[i].IsBinary || repoTreeEntries[i].OverSize || repoTreeEntries[i].OverPageSize {
			FilesNotHandled = append(FilesNotHandled, true)
		} else {
			FilesNotHandled = append(FilesNotHandled, false)
		}
	}

	ctx.Data["FilesContent"] = FilesContent
	ctx.Data["FilesFilename"] = FilesFilename
	ctx.Data["FilesNotHandled"] = FilesNotHandled

	ctx.Success(tmplEdit)
}

// EditPost POST
func EditPost(ctx *context.Context, f form.GitxtEdit) {
	if !ctx.IsLogged {
		ctx.Redirect(setting.AppSubURL + "/")
		return
	}

	if ctx.Data["LoggedUserID"] != ctx.Gitxt.Gitxt.UserID {
		ctx.Flash.Error(ctx.Tr("user.unauthorized"))
		ctx.Redirect(setting.AppSubURL + "/" + ctx.RepoOwnerUsername + "/" + ctx.Gitxt.Gitxt.Hash)
		return
	}

	ctx.Title("gitxt_edit.title")
	ctx.PageIs("GitxtEditPost")

	for i := range f.FilesFilename {
		// For each filename sanitize it
		f.FilesFilename[i] = sanitize.Filename(f.FilesFilename[i])
		if len(f.FilesFilename[i]) == 0 || f.FilesFilename[i] == "." {
			// If length is zero, use default filename
			f.FilesFilename[i] = fmt.Sprintf("gitxt%d.txt", i)
		}

		// No extension, forces .txt
		if len(filepath.Ext(f.FilesFilename[i])) == 0 {
			f.FilesFilename[i] = fmt.Sprintf("%s.txt", f.FilesFilename[i])
		}

		// Check if empty content
		if len(strings.TrimSpace(f.FilesContent[i])) <= 0 && !f.FilesNotHandled[i] {
			ctx.Data[fmt.Sprintf("Err_FilesContent_%d", i)] = ctx.Tr("gitxt_new.error_files_content")
			ctx.Data["HasError"] = true
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_new.files_content_cannot_empty")
			ctx.Flash.Error(ctx.Tr("gitxt_new.files_content_cannot_empty"), true)
		}
	}

	// Since the validation of slices doesn't works in bindings, manually update context
	ctx.Data["FilesFilename"] = f.FilesFilename
	ctx.Data["FilesContent"] = f.FilesContent
	ctx.Data["FilesNotHandled"] = f.FilesNotHandled
	ctx.Data["ExpiryHours"] = f.ExpiryHours

	// We got an error in the manual validation step, render with error
	if ctx.HasError() {
		ctx.Success(tmplNew)
		return
	}

	// Ok we are good

	repositoryUser := ctx.Gitxt.User.UserName
	repositoryHash := ctx.Gitxt.Gitxt.Hash

	repoPath := repository.RepoPath(repositoryUser, repositoryHash)

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Warn("Could not open repository %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.could_not_open"), true)
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	// Test if repository is empty
	isEmpty, err := repo.IsEmpty()
	if err != nil || isEmpty {
		log.Warn("Empty repository or corrupted %s: %s", ctx.Gitxt.Gitxt.Hash, err)
		ctx.Flash.Error(ctx.Tr("gitxt_git.repo_corrupt_or_empty"), true)
		ctx.Handle(500, "GitxtEditPost", err)
		return
	}

	// git whatever create blob
	// If it's an unhandled file : we get the Oid from filename then use it as-is
	// Else we create a blob and get his Oid
	var blobs []*git.Oid
	for i := range f.FilesFilename {
		var blob *git.Oid

		if f.FilesNotHandled[i] {
			blob, _ = gite.GetTreeFileOid(repo, f.FilesFilename[i])
		} else {
			blob, err = repo.CreateBlobFromBuffer([]byte(f.FilesContent[i]))
			if err != nil {
				log.Warn("init_error_create_blob: %s", err)
				ctx.Data["HasError"] = true
				ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_create_blob")
				ctx.Flash.Error(ctx.Tr("gitxt_git.error_create_blob"), true)
				ctx.Success(tmplNew)
				return
			}
		}
		blobs = append(blobs, blob)
	}

	//
	repoIndex, err := repo.Index()
	if err != nil {
		log.Warn("init_error_get_index: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_index")
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_get_index"), true)
		ctx.Success(tmplNew)
		return
	}

	// Add the blobs to the index
	// git update-index --add --cacheinfo 100644 "$BLOB_ID" "myfile.txt"
	for i := range f.FilesFilename {
		indexEntry := &git.IndexEntry{
			Path: f.FilesFilename[i],
			Mode: git.FilemodeBlob,
			Id:   blobs[i],
		}
		log.Trace("Add blob %s fname %s", indexEntry.Id.String(), indexEntry.Path)

		if err = repoIndex.Add(indexEntry); err != nil {
			log.Warn("init_error_add_entry: %s", err)
			ctx.Data["HasError"] = true
			ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_add_entry")
			ctx.Flash.Error(ctx.Tr("gitxt_git.error_add_entry"), true)
			ctx.Success(tmplNew)
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
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_index_write_tree"), true)
		ctx.Success(tmplNew)
		return

	}

	// Get latest tree
	repoTree, err := repo.LookupTree(repoTreeOid)
	if err != nil {
		log.Warn("init_error_lookup_tree: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_lookup_tree")
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_lookup_tree"), true)
		ctx.Success(tmplNew)
		return

	}

	// Get repo head
	repoHead, err := repo.Head()
	if err != nil {
		log.Warn("git_error_get_head: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_head")
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_get_head"), true)
		ctx.Success(tmplNew)
		return
	}

	// Get latest commit
	headCommit, err := repo.LookupCommit(repoHead.Target())
	if err != nil {
		log.Warn("git_error_get_head_commit: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_get_head_commit")
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_get_head_commit"), true)
		ctx.Success(tmplNew)
		return
	}

	// NEW_COMMIT=$(echo "My commit message" | git commit-tree "$TREE_ID" -p "$PARENT_COMMIT")
	ruAuthor := &git.Signature{
		Name:  repositoryUser,
		Email: "autocommit@git.txt",
	}
	_, err = repo.CreateCommit("HEAD", ruAuthor, ruAuthor, "Autocommit from git.txt", repoTree, headCommit)
	if err != nil {
		log.Warn("init_error_commit: %s", err)
		ctx.Data["HasError"] = true
		ctx.Data["ErrorMsg"] = ctx.Tr("gitxt_git.error_commit")
		ctx.Flash.Error(ctx.Tr("gitxt_git.error_commit"), true)
		ctx.Success(tmplNew)
		return

	}

	// git update-ref "refs/heads/$MY_BRANCH" "$NEW_COMMIT" "$PARENT_COMMIT"

	// 4. Insert info in database
	ctx.Gitxt.Gitxt.Description = f.Description
	ctx.Gitxt.Gitxt.ExpiryHours = f.ExpiryHours
	ctx.Gitxt.Gitxt.ExpiryUnix = time.Now().Add(time.Hour * time.Duration(f.ExpiryHours)).Unix()

	if err := models.UpdateGitxt(ctx.Gitxt.Gitxt); err != nil {
		switch {
		default:
			ctx.Handle(500, "EditPost", err)
		}
		return
	}

	// 5. Return render to gitxt view page

	log.Trace("Edit Pushed repository %s - %d", ctx.Gitxt.Gitxt.Hash, ctx.Gitxt.Gitxt.ID)
	ctx.Redirect(setting.AppSubURL + "/" + repositoryUser + "/" + ctx.Gitxt.Gitxt.Hash)

}
