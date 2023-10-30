package component

import "github.com/xackery/wlk/walk"

type FileViewEntry struct {
	Icon    *walk.Icon
	Name    string
	Ext     string
	Size    string
	RawSize int
	checked bool
}
