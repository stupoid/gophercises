package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

func RemoveIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func SameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func getDbPath() (string, error) {
	var path string
	home, err := homedir.Dir()
	if err != nil {
		return path, err
	}
	path = filepath.Join(home, dbFile)
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return path, err
	}
	return path, nil
}
