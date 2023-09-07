//go:build windows

package fileinfo

import (
	"errors"
	"io/fs"
)

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
