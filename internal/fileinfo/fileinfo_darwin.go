//go:build darwin

package fileinfo

import (
	"errors"
	"fmt"
	"io/fs"
	"os/user"
	"strconv"
	"syscall"
)

var (
	ErrNoUser  = errors.New("no user information available")
	ErrNoGroup = errors.New("no group information available")
)

func UserName(info fs.FileInfo) (string, error) {
	stat, ok := stat(info)
	if !ok {
		return "", ErrNoUser
	}

	usr, err := user.LookupId(strconv.FormatUint(uint64(stat.Uid), 10))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNoUser, err)
	}
	return usr.Username, nil
}

func GroupName(info fs.FileInfo) (string, error) {
	stat, ok := stat(info)
	if !ok {
		return "", ErrNoGroup
	}

	grp, err := user.LookupGroupId(strconv.FormatUint(uint64(stat.Gid), 10))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNoGroup, err)
	}
	return grp.Name, nil
}

func stat(info fs.FileInfo) (*syscall.Stat_t, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return stat, ok
}
