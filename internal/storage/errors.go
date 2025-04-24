package storage

import "fmt"

var (
	ErrWhenOpenFile      = func(loc string) error { return fmt.Errorf("error when open the file at: %s", loc) }
	ErrCannotGetFileInfo = func(loc string) error { return fmt.Errorf("cannot get file info at: %s", loc) }
	ErrCannotReadName    = func(loc string) error { return fmt.Errorf("cannot read file at: %s", loc) }
)
