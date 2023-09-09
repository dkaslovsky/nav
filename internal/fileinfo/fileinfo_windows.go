//go:build windows

package fileinfo

import (
	"errors"
	"io/fs"
)

// TODO: consider using https://github.com/itchio/ox/blob/12c6ca18d236/winox/permissions_windows.go#L375

var (
	ErrNoUser  = errors.New("no user information available")
	ErrNoGroup = errors.New("no group information available")
)

func UserName(info fs.FileInfo) (string, error) {
	return "", ErrNoUser
}

func Group(info fs.FileInfo) (string, error) {
	return "", ErrNoGroup
}
