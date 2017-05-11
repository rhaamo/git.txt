package cron

import (
	log "gopkg.in/clog.v1"
	"github.com/gogits/cron"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"dev.sigpipe.me/dashie/git.txt/models"
	"time"
)

var c = cron.New()

func NewContext() {
	var (
		entry	*cron.Entry
		err	error
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
	c.Start()
}

// ListTasks returns all running cron tasks.
func ListTasks() []*cron.Entry {
	return c.Entries()
}