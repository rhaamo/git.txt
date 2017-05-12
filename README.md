git.txt
=======

# What

One day this will do something like Gist with pastes backed by git.

Right now it's a mess of crap.

# Current Features list
- User or Anonymous Gitxts
- Archive download in Zip or TarGz of Gitxts repositories
- Git pull over HTTP/S for all Gitxts
- Git push over HTTP/S for user Gitxts (Using Basic Auth)
- Text Highlighting with Highlight.JS
- Size limit per individual file
- Size limit for the whole page (only text files not over-size are counted)

# Planned Feature List
- Pull or Push over SSH for user Gitxts
- Display PDF, images, Markdown, rich texts
- Raw files of Gitxts available
- API for third party tools

# Use or build

    Install or build libgit2 0.25.x ONLY. NO Version less than 0.25 (hello Debian Stable)
    go get -u github.com/jteeuwen/go-bindata/...
    go get -v -insecure -u dev.sigpipe.me/dashie/git.txt

# Sources used

I learned playing with Macaron/Xorm etc. from Gogs sources so lot of logic have been reused from Gogs.

# License

MIT, Dashie for git.txt and Gogs contributors for reused Gogs parts.
