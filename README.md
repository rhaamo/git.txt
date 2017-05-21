git.txt
=======

# What

It's a Pastebin where all pastes are backed in is own Git repository.

# Current Features list
- User or Anonymous Gitxts
- Archive download in Zip or TarGz of Gitxts repositories
- Git pull over HTTP/S for all Gitxts
- Git push over HTTP/S for user Gitxts (Using Basic Auth)
- Text Highlighting with Highlight.JS
- Size limit per individual file
- Size limit for the whole page (only text files not over-size are counted)
- Viewing RAW content of a file or download if binary type, a RAW Size Limit apply to whatever blob is wanted
- Rendering of special types: Markdown, Images, PDF (via PDF.js)
- Line Numbers for text files

# Planned Feature List
- Pull or Push over SSH for user Gitxts
- Display other rich texts
- API for third party tools
- More tests coverage
- Paste expiry (1d, 2d, 3d, 4d, 5d, 6d, 7d, 1w, 1m, 1y, never)

# Current issues or misbehaviour
- If a Git.txt upload is pushed with binary contents (or maybe big files, even text) they are not handled using web interface and cannot be edited that way, it's still possible to use Git itself and clone/push.

# Use or build

    Install or build libgit2 0.25.x ONLY. NO Version less than 0.25 (hello Debian Stable)
    go get -u github.com/jteeuwen/go-bindata/...
    go get -v -insecure -u dev.sigpipe.me/dashie/git.txt

# Contact, issues
- Main contact: Dashie <dashie (at) sigpipe (dot) me>
- Main repository: https://dev.sigpipe.me/dashie/git.txt
- Main issue tracker: https://dev.sigpipe.me/dashie/git.txt/issues

# Sources used

I learned playing with Macaron/Xorm etc. from Gogs sources so lot of logic have been reused from Gogs.

# License

MIT, Dashie for git.txt and Gogs contributors for reused Gogs parts.
