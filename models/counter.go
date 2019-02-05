package models

import (
	"fmt"
)

// Counter struct
type Counter struct {
	ID    int64  `xorm:"pk autoincr"`
	Name  string `xorm:"UNIQUE NOT NULL"`
	Count int64
}

// Don't change them !
const (
	currentGitxts = "current_gitxts"
	managedGitxts = "total_gitxts"
)

/* Counter for Gitxts [current number] */

// GetCounterGitxts get it
func GetCounterGitxts() (*Counter, error) {
	c := &Counter{Name: currentGitxts}
	has, err := x.Get(c)
	if err != nil || !has {
		c.Count = 0
		return c, fmt.Errorf("no counter available")
	}
	return c, nil
}

// Update a counter
func updateCounterGitxts(e Engine, count int64) error {
	c := &Counter{
		Name:  currentGitxts,
		Count: count,
	}
	// Try to get the Counter and create if not existent
	cc, err := GetCounterGitxts()
	if cc.Count <= 0 && err != nil {
		// create
		sess := x.NewSession()
		defer sess.Close()
		if err = sess.Begin(); err != nil {
			return err
		}

		if _, err = sess.Insert(c); err != nil {
			return err
		}

		return sess.Commit()
	}

	cc.Count = count

	// Update it
	_, err = e.Id(cc.ID).AllCols().Update(cc)
	return err
}

// UpdateCounterGitxts with count
func UpdateCounterGitxts(u int64) error {
	return updateCounterGitxts(x, u)
}

/* Counter for Gitxts [total managed number] */

// GetCounterGitxtsManaged get it
func GetCounterGitxtsManaged() (*Counter, error) {
	c := &Counter{Name: managedGitxts}
	has, err := x.Get(c)
	if err != nil || !has {
		c.Count = 0
		return c, fmt.Errorf("no counter available")
	}
	return c, nil
}

// Update a counter
func updateCounterGitxtsManaged(e Engine, count int64) error {
	c := &Counter{
		Name:  managedGitxts,
		Count: count,
	}
	// Try to get the Counter and create if not existent
	cc, err := GetCounterGitxtsManaged()
	if cc.Count <= 0 && err != nil {
		// create
		sess := x.NewSession()
		defer sess.Close()
		if err = sess.Begin(); err != nil {
			return err
		}

		if _, err = sess.Insert(c); err != nil {
			return err
		}

		return sess.Commit()
	}

	cc.Count = count

	// Update it
	_, err = e.Id(cc.ID).AllCols().Update(cc)
	return err
}

// UpdateCounterGitxtsManaged with count
func UpdateCounterGitxtsManaged(c int64) error {
	return updateCounterGitxtsManaged(x, c)
}
