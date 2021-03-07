package main

import (
	//"context"
	//"log"
	//"os"
	"sync"
	//"sync/atomic"
	//"time"

	//"bazil.org/fuse"
	//"bazil.org/fuse/fs"
	//"github.com/diamondburned/arikawa/v2/discord"
	//"github.com/diamondburned/arikawa/v2/state"
  //"github.com/diamondburned/ningen"
	//"github.com/pkg/errors"
)

type DMS struct {
	FS    *Filesystem
	Inode uint64

	Name string

	mu       sync.Mutex
}
