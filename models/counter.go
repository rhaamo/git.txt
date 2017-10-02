package models

import (
	"fmt"
)

// Counter struct
type Counter struct {
	ID		int64	`xorm:"pk autoincr"`
	Name	string	`xorm:"UNIQUE NOT NULL"`
	Count	int64
}

// Don't change them !
const (
	current_gitxts = "current_gitxts"
	managed_gitxts = "total_gitxts"
)

/* Counter for Gitxts [current number] */
func GetCounterGitxts() (*Counter, error) {
	c := &Counter{Name: current_gitxts}
	has, err := x.Get(c)
	if err != nil || !has {
		c.Count = 0
		return c, fmt.Errorf("no counter available")
	}
	return c, nil
}

// Update an user
func updateCounterGitxts(e Engine, count int64) error {
	c := &Counter{
		Name: current_gitxts,
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

// UpdateUser with datas
func UpdateCounterGitxts(u int64) error {
	return updateCounterGitxts(x, u)
}

/* Counter for Gitxts [total managed number] */
func GetCounterGitxtsManaged() (*Counter, error) {
	c := &Counter{Name: managed_gitxts}
	has, err := x.Get(c)
	if err != nil || !has {
		c.Count = 0
		return c, fmt.Errorf("no counter available")
	}
	return c, nil
}

// Update an user
func updateCounterGitxtsManaged(e Engine, count int64) error {
	c := &Counter{
		Name: managed_gitxts,
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

// UpdateUser with datas
func UpdateCounterGitxtsManaged(c int64) error {
	return updateCounterGitxtsManaged(x, c)
}