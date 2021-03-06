package cron

import (
	"dev.sigpipe.me/dashie/git.txt/models"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"github.com/gogits/cron"
	log "gopkg.in/clog.v1"
	"time"
)

var c = cron.New()

// NewContext initialize Cron stuff
func NewContext() {
	var (
		entry *cron.Entry
		err   error
	)

	if setting.Cron.RepoArchiveCleanup.Enabled {
		log.Trace("Enabling RepoArchiveCleanup")
		entry, err = c.AddFunc("Repository archive cleanup", setting.Cron.RepoArchiveCleanup.Schedule, models.DeleteOldRepositoryArchives)
		if err != nil {
			log.Fatal(2, "Cron.(repository archive cleanup): %v", err)
		}
		if setting.Cron.RepoArchiveCleanup.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go models.DeleteOldRepositoryArchives()
		}
	}

	entry, err = c.AddFunc("Delete expired repositories", "@every 1h", models.DeleteExpiredRepositories)
	if err != nil {
		log.Fatal(2, "Cron.(delete expired repositories): %v", err)
	}
	entry.Next = time.Now().Add(time.Minute * 1)

	c.Start()
}

// ListTasks returns all running cron tasks.
func ListTasks() []*cron.Entry {
	return c.Entries()
}
