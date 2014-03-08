package gollector_metrics

/*
// int statfs(const char *path, struct statfs *buf);

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
	"math"
	"unsafe"
)

/*
For a given mountpoint `path`, returns an [4]interface{} tuple containing:

* The free storage on the disk in bytes - this includes root reserved storage.

* The available storage in `path` -- this does not include root's storage.

* The total number of space in `path`.

* True/False based on readonly status for the mount point.

Supplying a directory that is not a mount point results in undefined behavior.
*/
func FSUsage(path string) interface{} {
	cPath := C.CString(path)
	stat := C.go_statvfs(cPath)
	readonly := C.go_fs_readonly(cPath)

	defer C.free(unsafe.Pointer(stat))
	defer C.free(unsafe.Pointer(cPath))

	blocks := uint64(stat.f_blocks)
	avail := uint64(stat.f_bavail)

	free := float64(0)

	if avail != 0 {
		free = math.Ceil(((float64(blocks) - float64(avail)) / float64(blocks)) * 100)
	}

	return [4]interface{}{
		free,
		avail * uint64(stat.f_frsize),
		blocks * uint64(stat.f_frsize),
		readonly == 1,
	}
}
