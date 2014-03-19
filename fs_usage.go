package gollector_metrics

/*
#include <sys/statvfs.h>
#include <stdlib.h>
#include <assert.h>

struct statvfs* go_statvfs(const char *path) {
  struct statvfs *fsinfo;
  fsinfo = malloc(sizeof(struct statvfs));
  assert(fsinfo != NULL);
  statvfs(path, fsinfo);
  return fsinfo;
}

int go_fs_readonly(const char *path) {
  struct statvfs *fsinfo = go_statvfs(path);

  return (fsinfo->f_flag & ST_RDONLY) == ST_RDONLY;
}
*/
import "C"

import (
	"unsafe"
)

/*
Type returned by FSUsage
*/
type FSInfo struct {
	Free     uint64 // The free storage on the disk in megabytes - this includes root reserved storage.
	Avail    uint64 // The available storage in `path` -- this does not include root's storage.
	Blocks   uint64 // The total number of space in `path`.
	ReadOnly bool   // True/False based on readonly status for the mount point.
}

/*
For a given mountpoint `path`, returns an FSInfo struct. Supplying a directory
that is not a mount point results in undefined behavior.
*/
func FSUsage(path string) FSInfo {
	cPath := C.CString(path)
	stat := C.go_statvfs(cPath)
	readonly := C.go_fs_readonly(cPath)
	frsize := uint64(stat.f_frsize)

	defer C.free(unsafe.Pointer(stat))
	defer C.free(unsafe.Pointer(cPath))

	blocks := uint64(stat.f_blocks)
	avail := uint64(stat.f_bavail)
	free := uint64(stat.f_bfree)

	return FSInfo{
		Free:     free * frsize,
		Avail:    avail * frsize,
		Blocks:   blocks * frsize,
		ReadOnly: readonly == 1,
	}
}
