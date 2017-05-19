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
- Rendering of special types: Markdown, Images

# Planned Feature List
- Pull or Push over SSH for user Gitxts
- Display PDF, other rich texts
- API for third party tools

# Use or build

    Install or build libgit2 0.25.x ONLY. NO Version less than 0.25 (hello Debian Stable)
    go get -u github.com/jteeuwen/go-bindata/...
    go get -v -insecure -u dev.sigpipe.me/dashie/git.txt

# Sources used

I learned playing with Macaron/Xorm etc. from Gogs sources so lot of logic have been reused from Gogs.

# License

MIT, Dashie for git.txt and Gogs contributors for reused Gogs parts.
